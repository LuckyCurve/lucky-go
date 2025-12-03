package forex

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"

	"lucky-go/notify"
)

var (
	amount float64
	push   bool
)

// forexCmd 表示汇率查询命令
var forexCmd = &cobra.Command{
	Use:   "forex [from] [to]",
	Short: "查询货币汇率",
	Long: `从 Frankfurter API 查询实时汇率。

示例:
  lucky-go forex USD CNY              # 查询 1 USD = ? CNY
  lucky-go forex EUR USD --amount 100 # 查询 100 EUR = ? USD
  lucky-go forex USD CNY --push       # 查询并推送到 Telegram

支持的货币: USD, EUR, CNY, JPY, GBP, AUD, CAD, CHF, HKD, SGD 等`,
	Args: cobra.ExactArgs(2),
	RunE: runForex,
}

func init() {
	forexCmd.Flags().Float64VarP(&amount, "amount", "a", 1, "兑换金额")
	forexCmd.Flags().BoolVarP(&push, "push", "p", false, "推送结果到 Telegram")
}

// runForex 执行汇率查询
func runForex(cmd *cobra.Command, args []string) error {
	from := strings.ToUpper(args[0])
	to := strings.ToUpper(args[1])

	// 获取汇率
	result, err := GetExchangeRate(from, to, amount)
	if err != nil {
		return err
	}

	// 渲染表格
	renderForexTable(result)

	// 如果需要推送到 Telegram
	if push {
		message := formatTelegramMessage(result)
		if err := notify.SendTelegramMessage(message); err != nil {
			return fmt.Errorf("推送到 Telegram 失败: %w", err)
		}
		fmt.Println("\n成功推送汇率信息到 Telegram")
	}

	return nil
}

// renderForexTable 渲染汇率查询结果表格
func renderForexTable(result *ExchangeResult) {
	// 定义颜色函数
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()

	// 配置 Colorized 渲染器
	cfg := renderer.ColorizedConfig{
		Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
		Settings: tw.Settings{
			Separators: tw.Separators{BetweenColumns: tw.On, ShowHeader: tw.On},
			Lines:      tw.Lines{ShowTop: tw.On, ShowBottom: tw.On, ShowHeaderLine: tw.On},
		},
		Symbols: tw.NewSymbols(tw.StyleLight),
	}

	// 创建表格
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(cfg)),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
	)

	// 设置表头
	table.Header([]string{"", "💱 汇率查询结果"})

	// 添加数据行
	_ = table.Append([]string{"源货币", cyanBold(result.From)})
	_ = table.Append([]string{"目标货币", cyanBold(result.To)})
	_ = table.Append([]string{"汇率", greenBold(fmt.Sprintf("%.4f", result.Rate))})
	_ = table.Append([]string{"金额", fmt.Sprintf("%.2f %s = %.2f %s",
		result.Amount, result.From, result.Converted, result.To)})
	_ = table.Append([]string{"更新时间", result.UpdateDate})

	// 渲染表格
	_ = table.Render()
}

// formatTelegramMessage 格式化 Telegram 消息（使用纯英文标签和表格边框保持等宽对齐）
func formatTelegramMessage(result *ExchangeResult) string {
	// 预先格式化 Amount 行内容
	amountStr := fmt.Sprintf("%.2f %s = %.2f %s", result.Amount, result.From, result.Converted, result.To)

	return fmt.Sprintf(`💱 *Exchange Rate*
`+"```"+`
+--------+------------------------+
| From   | %-22s |
| To     | %-22s |
| Rate   | %-22.4f |
| Amount | %-22s |
| Date   | %-22s |
+--------+------------------------+
`+"```"+`
_Source: Frankfurter API_`,
		result.From,
		result.To,
		result.Rate,
		amountStr,
		result.UpdateDate,
	)
}

// NewCommand 为汇率模块创建并返回汇率查询命令。
func NewCommand() *cobra.Command {
	return forexCmd
}
