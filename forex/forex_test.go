package forex

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient 模拟 HTTP 客户端
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetExchangeRate(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建模拟 JSON 响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"amount":1,"base":"USD","date":"2025-12-03","rates":{"CNY":7.2456}}`)),
		}

		// 替换 HTTP 客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		result, err := GetExchangeRate("USD", "CNY", 1)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if result.From != "USD" {
			t.Errorf("expected From 'USD', got '%s'", result.From)
		}

		if result.To != "CNY" {
			t.Errorf("expected To 'CNY', got '%s'", result.To)
		}

		expectedRate := 7.2456
		if result.Rate != expectedRate {
			t.Errorf("expected rate %.4f, got %.4f", expectedRate, result.Rate)
		}

		if result.UpdateDate != "2025-12-03" {
			t.Errorf("expected date '2025-12-03', got '%s'", result.UpdateDate)
		}
	})

	t.Run("WithAmount", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建模拟 JSON 响应（100 USD = 724.56 CNY）
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"amount":100,"base":"USD","date":"2025-12-03","rates":{"CNY":724.56}}`)),
		}

		// 替换 HTTP 客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		result, err := GetExchangeRate("USD", "CNY", 100)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if result.Amount != 100 {
			t.Errorf("expected Amount 100, got %.2f", result.Amount)
		}

		expectedConverted := 724.56
		if result.Converted != expectedConverted {
			t.Errorf("expected Converted %.2f, got %.2f", expectedConverted, result.Converted)
		}

		// 汇率应该是 724.56 / 100 = 7.2456
		expectedRate := 7.2456
		if result.Rate != expectedRate {
			t.Errorf("expected rate %.4f, got %.4f", expectedRate, result.Rate)
		}
	})

	t.Run("APIError", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建错误响应
		mockResponse := &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(`{"message":"error"}`)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		_, err := GetExchangeRate("USD", "CNY", 1)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("CurrencyNotFound", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建不包含目标货币的响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"amount":1,"base":"USD","date":"2025-12-03","rates":{"EUR":0.92}}`)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		_, err := GetExchangeRate("USD", "CNY", 1)
		if err == nil {
			t.Error("expected error for currency not found, got nil")
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建无效 JSON 响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`invalid json`)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		_, err := GetExchangeRate("USD", "CNY", 1)
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestGetSupportedCurrencies(t *testing.T) {
	currencies := GetSupportedCurrencies()

	if len(currencies) == 0 {
		t.Error("expected non-empty currency list")
	}

	// 验证包含常用货币
	expectedCurrencies := []string{"USD", "EUR", "CNY", "JPY", "GBP"}
	for _, expected := range expectedCurrencies {
		found := false
		for _, c := range currencies {
			if c == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected currency %s not found in list", expected)
		}
	}
}

func TestForexCommand(t *testing.T) {
	t.Run("CommandStructure", func(t *testing.T) {
		cmd := NewCommand()

		if cmd.Use != "forex [from] [to]" {
			t.Errorf("expected command use 'forex [from] [to]', got '%s'", cmd.Use)
		}

		if cmd.Short != "查询货币汇率" {
			t.Errorf("expected different short description, got '%s'", cmd.Short)
		}
	})

	t.Run("Flags", func(t *testing.T) {
		cmd := NewCommand()

		// 检查 amount flag
		amountFlag := cmd.Flags().Lookup("amount")
		if amountFlag == nil {
			t.Error("expected 'amount' flag to exist")
		}

		// 检查 push flag
		pushFlag := cmd.Flags().Lookup("push")
		if pushFlag == nil {
			t.Error("expected 'push' flag to exist")
		}
	})
}
