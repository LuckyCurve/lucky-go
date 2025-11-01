// Package finance provides financial calculations and data retrieval functionality.
package finance

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// peCmd represents the price-to-earnings ratio calculation command
var peCmd = &cobra.Command{
	Use:   "pe",
	Short: "Calculate finance PE ratios based on treasury and AAA company yields",
	Long:  `Calculate price-to-earnings ratios using current 10-year treasury and AAA corporate bond yields as benchmarks.`,
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

		// 输出 10 年期国债收益率相关计算
		fmt.Printf("10 years treasury earning rate: %v\n", treasuryResult.value)
		fmt.Printf("50%% : %.2f\n", 50/treasuryResult.value)
		fmt.Printf("75%% : %.2f\n", 75/treasuryResult.value)
		fmt.Printf("100%% : %.2f\n", 100/treasuryResult.value)
		fmt.Printf("125%% : %.2f\n", 125/treasuryResult.value)
		fmt.Printf("150%% : %.2f\n", 150/treasuryResult.value)

		fmt.Print("\n======================\n\n")

		// 输出 AAA 公司债券收益率相关计算
		fmt.Printf("aaa company rate: %v\n", aaaResult.value)
		fmt.Printf("50%% : %.2f\n", 50/aaaResult.value)
		fmt.Printf("75%% : %.2f\n", 75/aaaResult.value)
		fmt.Printf("100%% : %.2f\n", 100/aaaResult.value)
		fmt.Printf("125%% : %.2f\n", 125/aaaResult.value)
		fmt.Printf("150%% : %.2f\n", 150/aaaResult.value)

		return nil
	},
}

// HTTPClient 定义了一个HTTP客户端接口，用于模拟HTTP请求
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// 默认HTTP客户端
var defaultHTTPClient HTTPClient = &http.Client{}

// get10YearTreasuryYield retrieves the current 10-year treasury yield from macromicro.me.
// It returns the yield value as a float64 and any error encountered during the process.
func get10YearTreasuryYield() (float64, error) {
	return getTreasuryYield("https://sc.macromicro.me/series/354/10year-bond-yield", "#panel > main > div.mm-chart-collection > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val")
}

// getAAACompanyYield retrieves the current AAA corporate bond yield from macromicro.me.
// It returns the yield value as a float64 and any error encountered during the process.
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

// NewCommand creates and returns the PE calculation command for the finance module.
func NewCommand() *cobra.Command {
	return peCmd
}
