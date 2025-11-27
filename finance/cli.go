// Package finance 提供金融计算和数据检索功能。
package finance

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// peCmd 表示市盈率计算命令
var peCmd = &cobra.Command{
	Use:   "pe",
	Short: "基于国债和AAA公司收益率计算金融市盈率",
	Long:  `使用当前10年期国债和AAA公司债券收益率作为基准计算市盈率。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 使用通道接收结果和错误
		type result struct {
			value float64
			err   error
		}

		treasuryCh := make(chan result, 1)
		aaaCh := make(chan result, 1)

		// 并行获取 10 年期国债收益率
		go func() {
			value, err := get10YearTreasuryYield()
			treasuryCh <- result{value: value, err: err}
		}()

		// 并行获取 AAA 公司债券收益率
		go func() {
			value, err := getAAACompanyYield()
			aaaCh <- result{value: value, err: err}
		}()

		// 等待两个请求完成
		treasuryResult := <-treasuryCh
		if treasuryResult.err != nil {
			return treasuryResult.err
		}

		aaaResult := <-aaaCh
		if aaaResult.err != nil {
			return aaaResult.err
		}

		// 使用颜色定义
		titleColor := color.New(color.FgCyan, color.Bold)
		valueColor := color.New(color.FgBlue)
		pe50Color := color.New(color.FgGreen)            // 低估值
		pe75Color := color.New(color.FgYellow)           // 中等估值
		pe100Color := color.New(color.FgBlue)            // 基准估值
		pe125Color := color.New(color.FgRed)             // 高估值
		pe150Color := color.New(color.FgRed, color.Bold) // 很高估值

		// 输出 10 年期国债收益率相关计算
		titleColor.Println("┌─────────────────────┐")
		titleColor.Printf("│ 📊 国债收益率 PE 计算\n")
		fmt.Printf("│ 收益率: ")
		valueColor.Printf("%.2f%%\n", treasuryResult.value)
		titleColor.Println("├─────────────────────┤")

		fmt.Printf("│ ")
		pe50Color.Printf("%-8s %.2f\n", "50% PE:", 50/treasuryResult.value)
		fmt.Printf("│ ")
		pe75Color.Printf("%-8s %.2f\n", "75% PE:", 75/treasuryResult.value)
		fmt.Printf("│ ")
		pe100Color.Printf("%-8s %.2f\n", "100% PE:", 100/treasuryResult.value)
		fmt.Printf("│ ")
		pe125Color.Printf("%-8s %.2f\n", "125% PE:", 125/treasuryResult.value)
		fmt.Printf("│ ")
		pe150Color.Printf("%-8s %.2f\n", "150% PE:", 150/treasuryResult.value)
		titleColor.Println("└─────────────────────┘")

		fmt.Print("\n\n")

		// 输出 AAA 公司债券收益率相关计算
		titleColor.Println("┌─────────────────────┐")
		titleColor.Printf("│ 📊 AAA债券收益率 PE 计算\n")
		fmt.Printf("│ 收益率: ")
		valueColor.Printf("%.2f%%\n", aaaResult.value)
		titleColor.Println("├─────────────────────┤")

		fmt.Printf("│ ")
		pe50Color.Printf("%-8s %.2f\n", "50% PE:", 50/aaaResult.value)
		fmt.Printf("│ ")
		pe75Color.Printf("%-8s %.2f\n", "75% PE:", 75/aaaResult.value)
		fmt.Printf("│ ")
		pe100Color.Printf("%-8s %.2f\n", "100% PE:", 100/aaaResult.value)
		fmt.Printf("│ ")
		pe125Color.Printf("%-8s %.2f\n", "125% PE:", 125/aaaResult.value)
		fmt.Printf("│ ")
		pe150Color.Printf("%-8s %.2f\n", "150% PE:", 150/aaaResult.value)
		titleColor.Println("└─────────────────────┘")

		return nil
	},
}

// HTTPClient 定义了一个HTTP客户端接口，用于模拟HTTP请求
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// 默认HTTP客户端
var defaultHTTPClient HTTPClient = &http.Client{}

// get10YearTreasuryYield 从 macromicro.me 获取当前10年期国债收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func get10YearTreasuryYield() (float64, error) {
	return getTreasuryYield("https://sc.macromicro.me/series/354/10year-bond-yield", "#panel > main > div.mm-chart-collection > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val")
}

// getAAACompanyYield 从 macromicro.me 获取当前AAA公司债券收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func getAAACompanyYield() (float64, error) {
	return getTreasuryYield("https://sc.macromicro.me/series/618/moodys-aaa", "#panel > main > div > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val")
}

// getTreasuryYield 是一个通用函数，用于获取财务收益率数据
func getTreasuryYield(url, selector string) (float64, error) {
	// 构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	// 模拟 Chrome 浏览器 UA
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}

	// 用 CSS selector 抓取
	selection := doc.Find(selector).First()
	text := strings.TrimSpace(selection.Text())

	// 转 float
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// NewCommand 为金融模块创建并返回市盈率计算命令。
func NewCommand() *cobra.Command {
	return peCmd
}
