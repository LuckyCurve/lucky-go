// Package forex 提供汇率查询功能。
package forex

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// Frankfurter API 基础 URL（免费、无需 API Key）
	frankfurterAPIBaseURL = "https://api.frankfurter.app/latest"
)

// FrankfurterResponse 表示 Frankfurter API 的响应结构
type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// ExchangeResult 表示汇率查询结果
type ExchangeResult struct {
	From       string  // 源货币
	To         string  // 目标货币
	Rate       float64 // 汇率
	Amount     float64 // 源金额
	Converted  float64 // 转换后金额
	UpdateDate string  // 更新日期
}

// HTTPClient 定义了一个 HTTP 客户端接口，用于模拟 HTTP 请求
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// 默认 HTTP 客户端
var defaultHTTPClient HTTPClient = &http.Client{}

// GetExchangeRate 从 Frankfurter API 获取汇率
// from: 源货币代码 (如 USD)
// to: 目标货币代码 (如 CNY)
// amount: 兑换金额
func GetExchangeRate(from, to string, amount float64) (*ExchangeResult, error) {
	// 构建 API URL
	url := fmt.Sprintf("%s?from=%s&to=%s&amount=%.2f",
		frankfurterAPIBaseURL, from, to, amount)

	// 构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析 JSON 响应
	var apiResp FrankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 获取目标货币汇率
	rate, ok := apiResp.Rates[to]
	if !ok {
		return nil, fmt.Errorf("未找到货币 %s 的汇率", to)
	}

	// 计算实际汇率（API 返回的是转换后的金额，需要除以 amount 得到汇率）
	actualRate := rate / amount

	return &ExchangeResult{
		From:       from,
		To:         to,
		Rate:       actualRate,
		Amount:     amount,
		Converted:  rate,
		UpdateDate: apiResp.Date,
	}, nil
}

// GetSupportedCurrencies 返回常用的货币代码列表
func GetSupportedCurrencies() []string {
	return []string{
		"USD", "EUR", "CNY", "JPY", "GBP", "AUD", "CAD", "CHF", "HKD", "SGD",
	}
}
