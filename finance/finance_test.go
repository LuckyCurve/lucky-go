package finance

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// 创建一个模拟的HTTP客户端
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGet10YearTreasuryYield(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建模拟响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`<html><body><div id="panel"><main><div class="mm-chart-collection"><div class="mm-cc-hd"><div><div class="mm-cc-chart-stats-title pb-2 d-flex flex-wrap align-items-baseline"><div class="stat-val"><span class="val">3.50</span></div></div></div></div></main></div></body></html>`)),
		}

		// 替换HTTP客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		result, err := get10YearTreasuryYield()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 3.50
		if result != expected {
			t.Errorf("expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("InvalidResponse", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建错误响应
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, &url.Error{Op: "Get", URL: req.URL.String(), Err: errors.New("network error")}
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		_, err := get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("InvalidHTML", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建包含无效数据的响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`<html><body>Invalid HTML</body></html>`)),
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

		_, err := get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error for invalid HTML, got nil")
		}
	})

	t.Run("InvalidNumber", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建包含非数字文本的响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`<html><body><div id="panel"><main><div class="mm-chart-collection"><div class="mm-cc-hd"><div><div class="mm-cc-chart-stats-title pb-2 d-flex flex-wrap align-items-baseline"><div class="stat-val"><span class="val">not_a_number</span></div></div></div></div></main></div></body></html>`)),
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

		_, err := get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error for invalid number, got nil")
		}
	})
}

func TestGetAAACompanyYield(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 创建模拟响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`<html><body><div id="panel"><main><div><div class="mm-cc-hd"><div><div class="mm-cc-chart-stats-title pb-2 d-flex flex-wrap align-items-baseline"><div class="stat-val"><span class="val">4.25</span></div></div></div></div></main></div></body></html>`)),
		}

		// 替换HTTP客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端
		defer func() {
			defaultHTTPClient = originalClient
		}()

		result, err := getAAACompanyYield()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 4.25
		if result != expected {
			t.Errorf("expected %.2f, got %.2f", expected, result)
		}
	})
}

func TestPECommand(t *testing.T) {
	// 测试PE命令的结构，但不执行网络请求
	t.Run("CommandStructure", func(t *testing.T) {
		cmd := NewCommand()

		if cmd.Use != "pe" {
			t.Errorf("expected command use 'pe', got '%s'", cmd.Use)
		}

		if cmd.Short != "Calculate finance PE ratios based on treasury and AAA company yields" {
			t.Errorf("expected different short description, got '%s'", cmd.Short)
		}
	})
}
