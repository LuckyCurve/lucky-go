package finance

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
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
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建模拟 JSON 响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"observations":[{"date":"2024-01-01","value":"3.50"}]}`)),
		}

		// 替换HTTP客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		result, err := Get10YearTreasuryYield()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 3.50
		if result != expected {
			t.Errorf("expected %.2f, got %.2f", expected, result)
		}
	})

	t.Run("InvalidResponse", func(t *testing.T) {
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建错误响应
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, &url.Error{Op: "Get", URL: req.URL.String(), Err: errors.New("network error")}
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		_, err := Get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("EmptyObservations", func(t *testing.T) {
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建空观测数据的响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"observations":[]}`)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		_, err := Get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error for empty observations, got nil")
		}
	})

	t.Run("InvalidNumber", func(t *testing.T) {
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建包含非数字文本的响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"observations":[{"date":"2024-01-01","value":"not_a_number"}]}`)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		_, err := Get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error for invalid number, got nil")
		}
	})

	t.Run("MissingAPIKey", func(t *testing.T) {
		// 保存原始环境变量
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 清除 API Key
		os.Unsetenv("FRED_API_KEY")

		// 恢复环境变量
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		_, err := Get10YearTreasuryYield()
		if err == nil {
			t.Error("expected error for missing API key, got nil")
		}
	})
}

func TestGetAAACompanyYield(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建模拟 JSON 响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"observations":[{"date":"2024-01-01","value":"4.25"}]}`)),
		}

		// 替换HTTP客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		result, err := GetAAACompanyYield()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 4.25
		if result != expected {
			t.Errorf("expected %.2f, got %.2f", expected, result)
		}
	})
}

func TestGetBAAYield(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端和环境变量
		originalClient := defaultHTTPClient
		originalAPIKey := os.Getenv("FRED_API_KEY")

		// 设置测试环境变量
		os.Setenv("FRED_API_KEY", "test_api_key")

		// 创建模拟 JSON 响应
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"observations":[{"date":"2024-01-01","value":"5.50"}]}`)),
		}

		// 替换HTTP客户端
		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		// 恢复原始客户端和环境变量
		defer func() {
			defaultHTTPClient = originalClient
			if originalAPIKey == "" {
				os.Unsetenv("FRED_API_KEY")
			} else {
				os.Setenv("FRED_API_KEY", originalAPIKey)
			}
		}()

		result, err := GetBAAYield()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 5.50
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

		if cmd.Short != "基于国债和AAA公司收益率计算金融市盈率" {
			t.Errorf("expected different short description, got '%s'", cmd.Short)
		}
	})
}
