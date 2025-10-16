package finance

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "pe",
	Short: "calculate finance pe",
	RunE: func(cmd *cobra.Command, args []string) error {
		earningRate, err := get10YearTreasuryYield()
		if err != nil {
			return err
		}

		fmt.Printf("10 years treasury earning rate: %v\n", earningRate)
		fmt.Printf("50%% : %.2f\n", 50/earningRate)
		fmt.Printf("75%% : %.2f\n", 75/earningRate)
		fmt.Printf("100%% : %.2f\n", 100/earningRate)
		fmt.Printf("125%% : %.2f\n", 125/earningRate)
		fmt.Printf("150%% : %.2f\n", 150/earningRate)

		fmt.Print("\n======================\n\n")
		earningRate, err = getAAACompanyYield()
		if err != nil {
			return err
		}

		fmt.Printf("aaa company rate: %v\n", earningRate)
		fmt.Printf("50%% : %.2f\n", 50/earningRate)
		fmt.Printf("75%% : %.2f\n", 75/earningRate)
		fmt.Printf("100%% : %.2f\n", 100/earningRate)
		fmt.Printf("125%% : %.2f\n", 125/earningRate)
		fmt.Printf("150%% : %.2f\n", 150/earningRate)

		return nil
	},
}

func get10YearTreasuryYield() (float64, error) {
	// 构造请求
	req, err := http.NewRequest("GET", "https://sc.macromicro.me/series/354/10year-bond-yield", nil)
	if err != nil {
		return 0, err
	}
	// 模拟 Chrome 浏览器 UA
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
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
	selection := doc.Find("#panel > main > div.mm-chart-collection > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val").First()
	text := strings.TrimSpace(selection.Text())

	// 转 float
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func getAAACompanyYield() (float64, error) {
	// 构造请求
	req, err := http.NewRequest("GET", "https://sc.macromicro.me/series/618/moodys-aaa", nil)
	if err != nil {
		return 0, err
	}
	// 模拟 Chrome 浏览器 UA
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
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
	selection := doc.Find("#panel > main > div > div.mm-cc-hd > div > div.mm-cc-chart-stats-title.pb-2.d-flex.flex-wrap.align-items-baseline > div.stat-val > span.val").First()
	text := strings.TrimSpace(selection.Text())

	// 转 float
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func NewCommand() *cobra.Command {
	return sshCmd
}
