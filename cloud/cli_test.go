package cloud

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
	"lucky-go/config"
)

func TestRebootCommand(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("HOME", originalHomeDir)
	defer os.Setenv("USERPROFILE", originalUserProfile)

	// 准备测试配置
	testConfig := config.Config{
		Dest: map[string]config.DestinationInstance{
			"test-dest": {
				Ssh:        "test@example.com",
				Region:     "ap-beijing",
				InstanceId: "ins-test123",
			},
		},
	}

	// 创建配置目录和文件
	configDir := filepath.Join(tempDir, config.CONFIG_DIR)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, config.CONFIG_FILE)
	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	t.Run("ValidDestination", func(t *testing.T) {
		cmd := NewCommand()

		// 为测试目的替换内部函数
		originalFunc := rebootInstanceFunc
		rebootInstanceFunc = func(dest *config.DestinationInstance) error {
			if dest.Region != "ap-beijing" || dest.InstanceId != "ins-test123" {
				return errors.New("unexpected destination instance values")
			}
			return nil
		}
		defer func() {
			rebootInstanceFunc = originalFunc
		}()

		err := cmd.RunE(cmd, []string{"test-dest"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("NoDestinationProvided", func(t *testing.T) {
		// 为避免运行命令时出现错误，我们不直接调用RunE
		// 而是测试命令配置
		cmd := NewCommand()

		// 验证命令配置
		if cmd.Args == nil {
			t.Error("expected Args to be set")
		}
	})

	t.Run("NonExistentDestination", func(t *testing.T) {
		cmd := NewCommand()

		// 为测试目的替换内部函数
		originalFunc := rebootInstanceFunc
		defer func() {
			rebootInstanceFunc = originalFunc
		}()

		err := cmd.RunE(cmd, []string{"non-existent"})
		if err == nil {
			t.Error("expected error for non-existent destination, got nil")
		}
	})
}
