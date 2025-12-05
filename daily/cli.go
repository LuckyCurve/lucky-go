// Package daily 提供每日综合金融报告功能。
package daily

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"lucky-go/finance"
	"lucky-go/forex"
	"lucky-go/notify"
	"lucky-go/valuation"
)

var (
	push      bool
	forexFrom string
	forexTo   string
	forexAmt  float64
)

// dailyCmd 表示每日综合报告命令
var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "生成每日综合金融报告",
	Long: `获取 PE 估值、CAPE 估值和汇率数据，生成综合报告。

示例:
  lucky-go daily                              # 显示综合报告
  lucky-go daily --push                       # 推送到 Telegram
  lucky-go daily --forex-from USD --forex-to CNY  # 指定汇率货币对`,
	RunE: runDaily,
}

func init() {
	dailyCmd.Flags().BoolVarP(&push, "push", "p", false, "推送结果到 Telegram")
	dailyCmd.Flags().StringVar(&forexFrom, "forex-from", "USD", "汇率源货币")
	dailyCmd.Flags().StringVar(&forexTo, "forex-to", "CNY", "汇率目标货币")
	dailyCmd.Flags().Float64Var(&forexAmt, "forex-amount", 1, "汇率兑换金额")
}

// NewCommand 返回 daily 命令
func NewCommand() *cobra.Command {
	return dailyCmd
}

// DailyReport 包含每日报告所需的所有数据
type DailyReport struct {
	// PE 数据
	Treasury float64
	AAA      float64
	BAA      float64

	// CAPE 数据
	CAPE    float64
	FairPE  float64
	Premium float64

	// Forex 数据
	ForexResult *forex.ExchangeResult
}

func runDaily(cmd *cobra.Command, args []string) error {
	// 定义结果类型
	type floatResult struct {
		value float64
		err   error
	}
	type forexResult struct {
		value *forex.ExchangeResult
		err   error
	}

	// 创建通道
	treasuryCh := make(chan floatResult, 1)
	aaaCh := make(chan floatResult, 1)
	baaCh := make(chan floatResult, 1)
	capeCh := make(chan floatResult, 1)
	forexCh := make(chan forexResult, 1)

	// 并行获取所有数据
	go func() {
		value, err := finance.Get10YearTreasuryYield()
		treasuryCh <- floatResult{value: value, err: err}
	}()

	go func() {
		value, err := finance.GetAAACompanyYield()
		aaaCh <- floatResult{value: value, err: err}
	}()

	go func() {
		value, err := finance.GetBAAYield()
		baaCh <- floatResult{value: value, err: err}
	}()

	go func() {
		value, err := valuation.GetShillerCAPE()
		capeCh <- floatResult{value: value, err: err}
	}()

	go func() {
		result, err := forex.GetExchangeRate(forexFrom, forexTo, forexAmt)
		forexCh <- forexResult{value: result, err: err}
	}()

	// 收集结果
	report := &DailyReport{}

	treasuryRes := <-treasuryCh
	if treasuryRes.err != nil {
		return fmt.Errorf("获取国债收益率失败: %w", treasuryRes.err)
	}
	report.Treasury = treasuryRes.value

	aaaRes := <-aaaCh
	if aaaRes.err != nil {
		return fmt.Errorf("获取 AAA 收益率失败: %w", aaaRes.err)
	}
	report.AAA = aaaRes.value

	baaRes := <-baaCh
	if baaRes.err != nil {
		return fmt.Errorf("获取 BAA 收益率失败: %w", baaRes.err)
	}
	report.BAA = baaRes.value

	capeRes := <-capeCh
	if capeRes.err != nil {
		return fmt.Errorf("获取 CAPE 失败: %w", capeRes.err)
	}
	report.CAPE = capeRes.value
	report.FairPE = 100 / report.Treasury
	report.Premium = (report.CAPE - report.FairPE) / report.FairPE * 100

	forexRes := <-forexCh
	if forexRes.err != nil {
		return fmt.Errorf("获取汇率失败: %w", forexRes.err)
	}
	report.ForexResult = forexRes.value

	// 显示报告
	renderDailyReport(report)

	// 推送到 Telegram
	if push {
		message := formatDailyMessage(report)
		if err := notify.SendTelegramMessage(message); err != nil {
			return fmt.Errorf("推送到 Telegram 失败: %w", err)
		}
		fmt.Println("\n成功推送每日综合报告到 Telegram")
	}

	return nil
}

// formatDailyMessage 格式化每日综合报告为 Telegram 消息
func formatDailyMessage(r *DailyReport) string {
	// CAPE 评级
	var rating string
	if r.Premium > 50 {
		rating = "⚠️ 估值偏高"
	} else if r.Premium > 20 {
		rating = "📊 略高于合理"
	} else if r.Premium > -10 {
		rating = "✅ 合理区间"
	} else {
		rating = "🎯 估值偏低"
	}

	return fmt.Sprintf(`📰 *每日金融综合报告*
📅 %s

━━━━━━━━━━━━━━━━━━━━

📊 *PE 估值*

*国债基准 (%.2f%%)*
• 50%% PE: %.2f | 100%% PE: %.2f | 150%% PE: %.2f

*AAA 基准 (%.2f%%)*
• 50%% PE: %.2f | 100%% PE: %.2f | 150%% PE: %.2f

*BAA 基准 (%.2f%%)*
• 50%% PE: %.2f | 100%% PE: %.2f | 150%% PE: %.2f

━━━━━━━━━━━━━━━━━━━━

📈 *CAPE 估值*
• 席勒 CAPE: %.2f
• 合理 PE: %.2f
• 溢价/折价: %+.1f%%
• 评级: %s

━━━━━━━━━━━━━━━━━━━━

💱 *汇率信息*
• %s → %s
• 汇率: %.4f
• %.2f %s = %.2f %s

━━━━━━━━━━━━━━━━━━━━

_数据来源: FRED, Multpl, Frankfurter_`,
		time.Now().Format("2006-01-02"),
		// PE 数据
		r.Treasury,
		50/r.Treasury, 100/r.Treasury, 150/r.Treasury,
		r.AAA,
		50/r.AAA, 100/r.AAA, 150/r.AAA,
		r.BAA,
		50/r.BAA, 100/r.BAA, 150/r.BAA,
		// CAPE 数据
		r.CAPE, r.FairPE, r.Premium, rating,
		// Forex 数据
		r.ForexResult.From, r.ForexResult.To,
		r.ForexResult.Rate,
		r.ForexResult.Amount, r.ForexResult.From,
		r.ForexResult.Converted, r.ForexResult.To,
	)
}

// renderDailyReport 在终端渲染每日综合报告
func renderDailyReport(r *DailyReport) {
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()

	cfg := renderer.ColorizedConfig{
		Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
		Settings: tw.Settings{
			Separators: tw.Separators{BetweenColumns: tw.On, ShowHeader: tw.On},
			Lines:      tw.Lines{ShowTop: tw.On, ShowBottom: tw.On, ShowHeaderLine: tw.On},
		},
		Symbols: tw.NewSymbols(tw.StyleLight),
	}

	// PE 表格
	fmt.Println(cyanBold("\n📊 PE 估值"))
	peTable := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(cfg)),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
	)
	peTable.Header([]string{"", "国债", "AAA", "BAA"})
	_ = peTable.Append([]string{"收益率", fmt.Sprintf("%.2f%%", r.Treasury), fmt.Sprintf("%.2f%%", r.AAA), fmt.Sprintf("%.2f%%", r.BAA)})
	_ = peTable.Append([]string{"100% PE", greenBold(fmt.Sprintf("%.2f", 100/r.Treasury)), greenBold(fmt.Sprintf("%.2f", 100/r.AAA)), greenBold(fmt.Sprintf("%.2f", 100/r.BAA))})
	_ = peTable.Render()

	// CAPE 表格
	fmt.Println(cyanBold("\n📈 CAPE 估值"))
	var premiumColor func(a ...interface{}) string
	if r.Premium > 50 {
		premiumColor = redBold
	} else if r.Premium > 20 {
		premiumColor = yellowBold
	} else {
		premiumColor = greenBold
	}

	capeTable := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(cfg)),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
	)
	capeTable.Header([]string{"指标", "数值"})
	_ = capeTable.Append([]string{"席勒 CAPE", cyanBold(fmt.Sprintf("%.2f", r.CAPE))})
	_ = capeTable.Append([]string{"合理 PE", greenBold(fmt.Sprintf("%.2f", r.FairPE))})
	_ = capeTable.Append([]string{"溢价/折价", premiumColor(fmt.Sprintf("%+.1f%%", r.Premium))})
	_ = capeTable.Render()

	// Forex 表格
	fmt.Println(cyanBold("\n💱 汇率信息"))
	forexTable := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(cfg)),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
	)
	forexTable.Header([]string{"", "汇率查询"})
	_ = forexTable.Append([]string{"货币对", fmt.Sprintf("%s → %s", r.ForexResult.From, r.ForexResult.To)})
	_ = forexTable.Append([]string{"汇率", greenBold(fmt.Sprintf("%.4f", r.ForexResult.Rate))})
	_ = forexTable.Append([]string{"兑换", fmt.Sprintf("%.2f %s = %.2f %s", r.ForexResult.Amount, r.ForexResult.From, r.ForexResult.Converted, r.ForexResult.To)})
	_ = forexTable.Render()
}
