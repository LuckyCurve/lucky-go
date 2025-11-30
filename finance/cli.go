// Package finance 提供金融计算和数据检索功能。
package finance

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
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
		bbbCh := make(chan result, 1)

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

		// 并行获取 BBB 公司债券收益率
		go func() {
			value, err := getBBBYield()
			bbbCh <- result{value: value, err: err}
		}()

		// 等待三个请求完成
		treasuryResult := <-treasuryCh
		if treasuryResult.err != nil {
			return treasuryResult.err
		}

		aaaResult := <-aaaCh
		if aaaResult.err != nil {
			return aaaResult.err
		}

		bbbResult := <-bbbCh
		if bbbResult.err != nil {
			return bbbResult.err
		}

		// 使用 tablewriter 渲染合并的表格（三列并排）
		treasuryPEs := [5]float64{
			50 / treasuryResult.value,
			75 / treasuryResult.value,
			100 / treasuryResult.value,
			125 / treasuryResult.value,
			150 / treasuryResult.value,
		}
		aaaPEs := [5]float64{
			50 / aaaResult.value,
			75 / aaaResult.value,
			100 / aaaResult.value,
			125 / aaaResult.value,
			150 / aaaResult.value,
		}
		bbbPEs := [5]float64{
			50 / bbbResult.value,
			75 / bbbResult.value,
			100 / bbbResult.value,
			125 / bbbResult.value,
			150 / bbbResult.value,
		}

		renderThreeColumnPETable(
			"国债收益率", treasuryResult.value, treasuryPEs,
			"AAA债券收益率", aaaResult.value, aaaPEs,
			"BAA债券收益率", bbbResult.value, bbbPEs,
		)

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

// getBBBYield 从 macromicro.me 获取当前BBB/Baa公司债券收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func getBBBYield() (float64, error) {
	return getTreasuryYield("https://sc.macromicro.me/series/619/moodys-baa", "#panel > main > div > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val")
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

// renderThreeColumnPETable 渲染三列 PE 表格，包含国债、AAA和BBB债券数据
func renderThreeColumnPETable(title1 string, yield1 float64, pe1 [5]float64,
	title2 string, yield2 float64, pe2 [5]float64,
	title3 string, yield3 float64, pe3 [5]float64) {
	// 创建表格
	table := tablewriter.NewWriter(os.Stdout)

	// 设置表头（四列：标签、国债、AAA、BBB）
	table.SetHeader([]string{"", fmt.Sprintf("📊 %s", title1), fmt.Sprintf("📊 %s", title2), fmt.Sprintf("📊 %s", title3)})
	table.SetBorder(true)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetAutoWrapText(false)

	// 定义颜色样式
	lowPEColor := tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold}
	midPEColor := tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold}
	basePEColor := tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold}
	highPEColor := tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold}
	veryHighPEColor := tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold}

	// 添加数据行（四列：标签、国债值、AAA值、BBB值）
	labels := []string{"50% PE:", "75% PE:", "100% PE:", "125% PE:", "150% PE:"}
	colors := []tablewriter.Colors{lowPEColor, midPEColor, basePEColor, highPEColor, veryHighPEColor}

	for i := 0; i < 5; i++ {
		table.Rich([]string{
			labels[i],
			fmt.Sprintf("%.2f", pe1[i]),
			fmt.Sprintf("%.2f", pe2[i]),
			fmt.Sprintf("%.2f", pe3[i]),
		}, []tablewriter.Colors{
			{}, colors[i], colors[i], colors[i],
		})
	}

	// 添加收益率行
	table.Rich([]string{
		"收益率",
		fmt.Sprintf("%.2f%%", yield1),
		fmt.Sprintf("%.2f%%", yield2),
		fmt.Sprintf("%.2f%%", yield3),
	}, []tablewriter.Colors{
		{}, {}, {}, {},
	})

	// 渲染表格
	table.Render()
}

// NewCommand 为金融模块创建并返回市盈率计算命令。
func NewCommand() *cobra.Command {
	return peCmd
}
