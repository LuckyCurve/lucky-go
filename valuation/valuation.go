// Package valuation 提供市场估值数据获取功能。
package valuation

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const (
	// Multpl.com Shiller PE 页面 URL
	multplShillerPEURL = "https://www.multpl.com/shiller-pe"
)

// HTTPClient 定义 HTTP 客户端接口
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// 默认 HTTP 客户端
var defaultHTTPClient HTTPClient = &http.Client{}

// CAPEResult 表示 CAPE 查询结果
type CAPEResult struct {
	Value      float64 // 当前 CAPE 值
	Mean       float64 // 历史均值
	Median     float64 // 历史中位数
	UpdateDate string  // 更新日期
}

// GetShillerCAPE 从 Multpl.com 获取当前席勒 CAPE 值
func GetShillerCAPE() (float64, error) {
	req, err := http.NewRequest("GET", multplShillerPEURL, nil)
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 User-Agent 避免被拒绝
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; lucky-go/1.0)")

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	// 使用正则提取 CAPE 值
	// 匹配模式: "Current Shiller PE Ratio is XX.XX" 或数字
	re := regexp.MustCompile(`(?i)(?:current\s+)?shiller\s+pe\s+ratio\s+(?:is\s+)?(\d+\.?\d*)`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		// 备用模式：直接匹配页面中的大数字
		re = regexp.MustCompile(`>(\d{2}\.\d{2})<`)
		matches = re.FindSubmatch(body)
		if len(matches) < 2 {
			return 0, fmt.Errorf("无法从页面解析 CAPE 值")
		}
	}

	value, err := strconv.ParseFloat(string(matches[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("解析 CAPE 值失败: %w", err)
	}

	return value, nil
}
