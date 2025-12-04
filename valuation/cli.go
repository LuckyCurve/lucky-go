package valuation

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"lucky-go/finance"
	"lucky-go/notify"
)

var push bool

// capeCmd 表示 CAPE 估值命令
var capeCmd = &cobra.Command{
	Use:   "cape",
	Short: "查询标普500席勒CAPE估值",
	Long: `从 Multpl.com 获取当前席勒 CAPE (周期调整市盈率)，
并与基于国债收益率计算的合理 PE 进行对比。

示例:
  lucky-go cape         # 查询 CAPE 估值
  lucky-go cape --push  # 查询并推送到 Telegram`,
	RunE: runCAPE,
}

func init() {
	capeCmd.Flags().BoolVarP(&push, "push", "p", false, "推送结果到 Telegram")
}

// runCAPE 执行 CAPE 估值查询
func runCAPE(cmd *cobra.Command, args []string) error {
	// 并行获取数据
	type result struct {
		value float64
		err   error
	}

	capeCh := make(chan result, 1)
	treasuryCh := make(chan result, 1)

	go func() {
		value, err := GetShillerCAPE()
		capeCh <- result{value: value, err: err}
	}()

	go func() {
		value, err := finance.Get10YearTreasuryYield()
		treasuryCh <- result{value: value, err: err}
	}()

	// 等待结果
	capeResult := <-capeCh
	if capeResult.err != nil {
		return fmt.Errorf("获取 CAPE 失败: %w", capeResult.err)
	}

	treasuryResult := <-treasuryCh
	if treasuryResult.err != nil {
		return fmt.Errorf("获取国债收益率失败: %w", treasuryResult.err)
	}

	// 计算合理 PE (100% 档位)
	fairPE := 100 / treasuryResult.value

	// 计算溢价/折价
	premium := (capeResult.value - fairPE) / fairPE * 100

	// 渲染表格
	renderCAPETable(capeResult.value, fairPE, treasuryResult.value, premium)

	// 推送到 Telegram
	if push {
		message := formatCAPETelegramMessage(capeResult.value, fairPE, treasuryResult.value, premium)
		if err := notify.SendTelegramMessage(message); err != nil {
			return fmt.Errorf("推送到 Telegram 失败: %w", err)
		}
		fmt.Println("\n成功推送 CAPE 估值到 Telegram")
	}

	return nil
}

// renderCAPETable 渲染 CAPE 估值对比表格
func renderCAPETable(cape, fairPE, treasury, premium float64) {
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()

	// 根据溢价程度选择颜色
	var premiumColor func(a ...interface{}) string
	if premium > 50 {
		premiumColor = redBold
	} else if premium > 20 {
		premiumColor = yellowBold
	} else {
		premiumColor = greenBold
	}

	cfg := renderer.ColorizedConfig{
		Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
		Settings: tw.Settings{
			Separators: tw.Separators{BetweenColumns: tw.On, ShowHeader: tw.On},
			Lines:      tw.Lines{ShowTop: tw.On, ShowBottom: tw.On, ShowHeaderLine: tw.On},
		},
		Symbols: tw.NewSymbols(tw.StyleLight),
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(cfg)),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
	)

	table.Header([]string{"指标", "📊 市场估值对比"})

	_ = table.Append([]string{"席勒 CAPE", cyanBold(fmt.Sprintf("%.2f", cape))})
	_ = table.Append([]string{"合理 PE (国债基准)", greenBold(fmt.Sprintf("%.2f", fairPE))})
	_ = table.Append([]string{"10年期国债收益率", fmt.Sprintf("%.2f%%", treasury)})
	_ = table.Append([]string{"溢价/折价", premiumColor(fmt.Sprintf("%+.1f%%", premium))})

	_ = table.Render()

	// 打印估值评级
	fmt.Println()
	if premium > 50 {
		fmt.Println(redBold("⚠️  市场估值偏高，建议谨慎"))
	} else if premium > 20 {
		fmt.Println(yellowBold("📊 市场估值略高于合理水平"))
	} else if premium > -10 {
		fmt.Println(greenBold("✅ 市场估值处于合理区间"))
	} else {
		fmt.Println(greenBold("🎯 市场估值偏低，可能存在机会"))
	}
}

// formatCAPETelegramMessage 格式化 Telegram 消息
func formatCAPETelegramMessage(cape, fairPE, treasury, premium float64) string {
	var rating string
	if premium > 50 {
		rating = "⚠️ 估值偏高"
	} else if premium > 20 {
		rating = "📊 略高于合理"
	} else if premium > -10 {
		rating = "✅ 合理区间"
	} else {
		rating = "🎯 估值偏低"
	}

	return fmt.Sprintf(`📈 *标普500 CAPE 估值报告*

*市场估值*
• 席勒 CAPE: %.2f
• 合理 PE: %.2f
• 溢价/折价: %+.1f%%

*基准数据*
• 10年期国债: %.2f%%

*评级: %s*

_数据来源: Multpl.com, FRED_`,
		cape, fairPE, premium, treasury, rating)
}

// NewCommand 返回 CAPE 估值命令
func NewCommand() *cobra.Command {
	return capeCmd
}
