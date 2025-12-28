package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendTelegramMessage(t *testing.T) {
	tests := []struct {
		name        string
		botToken    string
		chatID      string
		message     string
		serverResp  TelegramResponse
		serverCode  int
		wantErr     bool
		errContains string
	}{
		{
			name:        "MissingBotToken",
			botToken:    "",
			chatID:      "123456",
			message:     "test message",
			wantErr:     true,
			errContains: "TELEGRAM_BOT_TOKEN",
		},
		{
			name:        "MissingChatID",
			botToken:    "test-token",
			chatID:      "",
			message:     "test message",
			wantErr:     true,
			errContains: "TELEGRAM_CHAT_ID",
		},
		{
			name:       "SuccessfulSend",
			botToken:   "test-token",
			chatID:     "123456",
			message:    "Hello, World!",
			serverResp: TelegramResponse{OK: true},
			serverCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "APIError",
			botToken:   "test-token",
			chatID:     "123456",
			message:    "test message",
			serverResp: TelegramResponse{OK: false, Description: "Bad Request: chat not found"},
			serverCode: http.StatusOK,
			wantErr:    true,
			errContains: "chat not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 保存原始环境变量
			origBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
			origChatID := os.Getenv("TELEGRAM_CHAT_ID")
			defer func() {
				os.Setenv("TELEGRAM_BOT_TOKEN", origBotToken)
				os.Setenv("TELEGRAM_CHAT_ID", origChatID)
			}()

			// 设置测试环境变量
			if tt.botToken != "" {
				os.Setenv("TELEGRAM_BOT_TOKEN", tt.botToken)
			} else {
				os.Unsetenv("TELEGRAM_BOT_TOKEN")
			}
			if tt.chatID != "" {
				os.Setenv("TELEGRAM_CHAT_ID", tt.chatID)
			} else {
				os.Unsetenv("TELEGRAM_CHAT_ID")
			}

			// 如果需要测试服务器响应，创建 mock server
			if tt.botToken != "" && tt.chatID != "" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// 验证请求
					if r.Method != http.MethodPost {
						t.Errorf("Expected POST request, got %s", r.Method)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
					}

					// 解析并验证请求体
					var msg TelegramMessage
					if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
						t.Errorf("Failed to decode request body: %v", err)
					}
					if msg.ChatID != tt.chatID {
						t.Errorf("Expected ChatID %s, got %s", tt.chatID, msg.ChatID)
					}
					if msg.Text != tt.message {
						t.Errorf("Expected Text %s, got %s", tt.message, msg.Text)
					}
					if msg.ParseMode != "Markdown" {
						t.Errorf("Expected ParseMode Markdown, got %s", msg.ParseMode)
					}

					// 返回响应
					w.WriteHeader(tt.serverCode)
					_ = json.NewEncoder(w).Encode(tt.serverResp)
				}))
				defer server.Close()

				// 替换 API URL
				origBaseURL := telegramAPIBaseURL
				setTelegramAPIBaseURL(server.URL + "/bot")
				defer setTelegramAPIBaseURL(origBaseURL)
			}

			err := SendTelegramMessage(tt.message)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Error %q should contain %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSendTelegramMessage_NetworkError(t *testing.T) {
	// 保存原始环境变量
	origBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	origChatID := os.Getenv("TELEGRAM_CHAT_ID")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origBotToken)
		os.Setenv("TELEGRAM_CHAT_ID", origChatID)
	}()

	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("TELEGRAM_CHAT_ID", "123456")

	// 使用无效的 URL 模拟网络错误
	origBaseURL := telegramAPIBaseURL
	setTelegramAPIBaseURL("http://invalid.invalid.invalid/bot")
	defer setTelegramAPIBaseURL(origBaseURL)

	err := SendTelegramMessage("test message")
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

func TestSendTelegramMessage_InvalidJSONResponse(t *testing.T) {
	// 保存原始环境变量
	origBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	origChatID := os.Getenv("TELEGRAM_CHAT_ID")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origBotToken)
		os.Setenv("TELEGRAM_CHAT_ID", origChatID)
	}()

	os.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	os.Setenv("TELEGRAM_CHAT_ID", "123456")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	origBaseURL := telegramAPIBaseURL
	setTelegramAPIBaseURL(server.URL + "/bot")
	defer setTelegramAPIBaseURL(origBaseURL)

	err := SendTelegramMessage("test message")
	if err == nil {
		t.Error("Expected JSON parse error, got nil")
	}
	if !containsString(err.Error(), "解析响应失败") {
		t.Errorf("Error should mention parse failure, got: %v", err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
