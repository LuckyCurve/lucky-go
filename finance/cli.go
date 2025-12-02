// Package finance 提供金融计算和数据检索功能。
package finance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

const (
	// FRED API 基础 URL
	fredAPIBaseURL = "https://api.stlouisfed.org/fred/series/observations"
	// FRED Series IDs
	seriesDGS10 = "DGS10" // 10年期国债收益率
	seriesAAA   = "AAA"   // AAA 公司债收益率
	seriesBAA   = "BAA"   // BAA 公司债收益率
)

// FredResponse 表示 FRED API 的响应结构
type FredResponse struct {
	Observations []FredObservation `json:"observations"`
}

// FredObservation 表示单个观测数据点
type FredObservation struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

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
			value, err := Get10YearTreasuryYield()
			treasuryCh <- result{value: value, err: err}
		}()

		// 并行获取 AAA 公司债券收益率
		go func() {
			value, err := GetAAACompanyYield()
			aaaCh <- result{value: value, err: err}
		}()

		// 并行获取 BBB 公司债券收益率
		go func() {
			value, err := GetBAAYield()
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

// Get10YearTreasuryYield 从 FRED API 获取当前10年期国债收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func Get10YearTreasuryYield() (float64, error) {
	return GetFredYield(seriesDGS10)
}

// GetAAACompanyYield 从 FRED API 获取当前AAA公司债券收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func GetAAACompanyYield() (float64, error) {
	return GetFredYield(seriesAAA)
}

// GetBAAYield 从 FRED API 获取当前BAA公司债券收益率。
// 它返回收益率值作为 float64 以及在此过程中遇到的任何错误。
func GetBAAYield() (float64, error) {
	return GetFredYield(seriesBAA)
}

// GetFredYield 从 FRED API 获取指定 series 的最新收益率数据
func GetFredYield(seriesID string) (float64, error) {
	// 获取 API Key
	apiKey := os.Getenv("FRED_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("FRED_API_KEY 环境变量未设置，请访问 https://fred.stlouisfed.org/docs/api/api_key.html 申请")
	}

	// 构建 API URL
	url := fmt.Sprintf("%s?series_id=%s&api_key=%s&file_type=json&sort_order=desc&limit=1",
		fredAPIBaseURL, seriesID, apiKey)

	// 构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("FRED API 请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析 JSON 响应
	var fredResp FredResponse
	if err := json.NewDecoder(resp.Body).Decode(&fredResp); err != nil {
		return 0, fmt.Errorf("解析 FRED API 响应失败: %w", err)
	}

	if len(fredResp.Observations) == 0 {
		return 0, fmt.Errorf("FRED API 未返回 %s 的数据", seriesID)
	}

	// 获取最新值并转换为 float64
	valueStr := fredResp.Observations[0].Value
	if valueStr == "." {
		return 0, fmt.Errorf("FRED API 返回的 %s 数据不可用", seriesID)
	}

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析收益率值失败: %w", err)
	}

	return val, nil
}

// renderThreeColumnPETable 渲染三列 PE 表格，包含国债、AAA和BBB债券数据
func renderThreeColumnPETable(title1 string, yield1 float64, pe1 [5]float64,
	title2 string, yield2 float64, pe2 [5]float64,
	title3 string, yield3 float64, pe3 [5]float64) {

	// 定义颜色函数
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	blueBold := color.New(color.FgBlue, color.Bold).SprintFunc()
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()

	colorFuncs := []func(a ...interface{}) string{greenBold, yellowBold, blueBold, redBold, redBold}

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
	table.Header([]string{"", fmt.Sprintf("📊 %s", title1), fmt.Sprintf("📊 %s", title2), fmt.Sprintf("📊 %s", title3)})

	// 添加数据行（手动着色值）
	labels := []string{"50% PE:", "75% PE:", "100% PE:", "125% PE:", "150% PE:"}

	for i := 0; i < 5; i++ {
		cf := colorFuncs[i]
		table.Append([]string{
			labels[i],
			cf(fmt.Sprintf("%.2f", pe1[i])),
			cf(fmt.Sprintf("%.2f", pe2[i])),
			cf(fmt.Sprintf("%.2f", pe3[i])),
		})
	}

	// 添加收益率行
	table.Append([]string{
		"收益率",
		fmt.Sprintf("%.2f%%", yield1),
		fmt.Sprintf("%.2f%%", yield2),
		fmt.Sprintf("%.2f%%", yield3),
	})

	// 渲染表格
	table.Render()
}

// NewCommand 为金融模块创建并返回市盈率计算命令。
func NewCommand() *cobra.Command {
	return peCmd
}
