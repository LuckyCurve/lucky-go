package valuation

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

func TestGetShillerCAPE(t *testing.T) {
	t.Run("ValidResponse", func(t *testing.T) {
		// 保存原始客户端
		originalClient := defaultHTTPClient

		// 模拟 Multpl.com 页面响应
		mockHTML := `<div id="current">Current Shiller PE Ratio is 40.45</div>`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(mockHTML)),
		}

		defaultHTTPClient = &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		}

		defer func() {
			defaultHTTPClient = originalClient
		}()

		result, err := GetShillerCAPE()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := 40.45
		if result != expected {
			t.Errorf("expected %.2f, got %.2f", expected, result)
		}
	})
}
