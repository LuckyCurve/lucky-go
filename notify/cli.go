// Package notify 提供通知推送功能。
package notify

import (
	"fmt"

	"github.com/spf13/cobra"

	"lucky-go/finance"
)

// pushCmd 表示推送命令
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "获取金融数据并推送到 Telegram",
	Long:  `从 FRED API 获取最新的收益率数据，计算 PE 值并推送到 Telegram。`,
	RunE:  runPush,
}

// runPush 执行数据获取和推送
func runPush(cmd *cobra.Command, args []string) error {
	// 并行获取收益率数据
	type result struct {
		value float64
		err   error
	}

	treasuryCh := make(chan result, 1)
	aaaCh := make(chan result, 1)
	baaCh := make(chan result, 1)

	go func() {
		value, err := finance.Get10YearTreasuryYield()
		treasuryCh <- result{value: value, err: err}
	}()

	go func() {
		value, err := finance.GetAAACompanyYield()
		aaaCh <- result{value: value, err: err}
	}()

	go func() {
		value, err := finance.GetBAAYield()
		baaCh <- result{value: value, err: err}
	}()

	// 等待结果
	treasuryResult := <-treasuryCh
	if treasuryResult.err != nil {
		return fmt.Errorf("获取国债收益率失败: %w", treasuryResult.err)
	}

	aaaResult := <-aaaCh
	if aaaResult.err != nil {
		return fmt.Errorf("获取 AAA 债券收益率失败: %w", aaaResult.err)
	}

	baaResult := <-baaCh
	if baaResult.err != nil {
		return fmt.Errorf("获取 BAA 债券收益率失败: %w", baaResult.err)
	}

	// 构建消息
	message := formatPEMessage(treasuryResult.value, aaaResult.value, baaResult.value)

	// 推送到 Telegram
	if err := SendTelegramMessage(message); err != nil {
		return fmt.Errorf("推送到 Telegram 失败: %w", err)
	}

	fmt.Println("成功推送 PE 数据到 Telegram")
	return nil
}

// formatPEMessage 格式化 PE 数据为 Telegram 消息
func formatPEMessage(treasury, aaa, baa float64) string {
	// 计算各档位 PE
	treasuryPE100 := 100 / treasury
	aaaPE100 := 100 / aaa
	baaPE100 := 100 / baa

	return fmt.Sprintf(`📊 *每日 PE 估值报告*

*收益率数据*
• 10年期国债: %.2f%%
• AAA 公司债: %.2f%%
• BAA 公司债: %.2f%%

*100%% PE 估值*
• 国债基准: %.2f
• AAA 基准: %.2f
• BAA 基准: %.2f

*PE 区间参考*
| 档位 | 国债 | AAA | BAA |
|------|------|-----|-----|
| 50%% | %.2f | %.2f | %.2f |
| 75%% | %.2f | %.2f | %.2f |
| 100%% | %.2f | %.2f | %.2f |
| 125%% | %.2f | %.2f | %.2f |
| 150%% | %.2f | %.2f | %.2f |

_数据来源: FRED (Federal Reserve Economic Data)_`,
		treasury, aaa, baa,
		treasuryPE100, aaaPE100, baaPE100,
		50/treasury, 50/aaa, 50/baa,
		75/treasury, 75/aaa, 75/baa,
		100/treasury, 100/aaa, 100/baa,
		125/treasury, 125/aaa, 125/baa,
		150/treasury, 150/aaa, 150/baa,
	)
}

// NewCommand 为通知模块创建并返回推送命令。
func NewCommand() *cobra.Command {
	return pushCmd
}
