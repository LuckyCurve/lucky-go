package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadDestinationInstance(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	// 在Windows上，UserHomeDir()使用USERPROFILE
	originalUserProfile := os.Getenv("USERPROFILE")
	
	// 设置临时目录为HOME和USERPROFILE
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	// 恢复原始值
	defer os.Setenv("HOME", originalHomeDir)
	defer os.Setenv("USERPROFILE", originalUserProfile)

	// 准备测试配置数据
	testConfig := Config{
		Dest: map[string]DestinationInstance{
			"test-dest": {
				Ssh:        "test@example.com",
				Region:     "ap-beijing",
				InstanceId: "ins-test123",
			},
		},
	}

	// 创建配置目录和文件
	configDir := filepath.Join(tempDir, CONFIG_DIR)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, CONFIG_FILE)
	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// 测试成功加载
	t.Run("LoadExistingDestination", func(t *testing.T) {
		instance, err := LoadDestinationInstance("test-dest")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if instance == nil {
			t.Fatal("expected instance, got nil")
		}
		if instance.Ssh != "test@example.com" {
			t.Errorf("expected SSH to be 'test@example.com', got '%s'", instance.Ssh)
		}
		if instance.Region != "ap-beijing" {
			t.Errorf("expected region to be 'ap-beijing', got '%s'", instance.Region)
		}
		if instance.InstanceId != "ins-test123" {
			t.Errorf("expected instance ID to be 'ins-test123', got '%s'", instance.InstanceId)
		}
	})

	// 测试加载不存在的实例
	t.Run("LoadNonExistentDestination", func(t *testing.T) {
		_, err := LoadDestinationInstance("non-existent")
		if err == nil {
			t.Error("expected error for non-existent destination, got nil")
		}
		expected := "not destination non-existent in config"
		if err.Error() != expected {
			t.Errorf("expected error '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestConfig_SaveConfig(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("HOME", originalHomeDir)
	defer os.Setenv("USERPROFILE", originalUserProfile)

	// 准备测试配置
	testConfig := Config{
		Dest: map[string]DestinationInstance{
			"save-test": {
				Ssh:        "save@example.com",
				Region:     "ap-shanghai",
				InstanceId: "ins-save456",
			},
		},
	}

	// 保存配置
	err := testConfig.SaveConfig()
	if err != nil {
		t.Fatalf("expected no error when saving config, got: %v", err)
	}

	// 验证配置文件是否正确创建
	configPath, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("failed to get config file path: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var loadedConfig Config
	err = yaml.Unmarshal(data, &loadedConfig)
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	// 验证配置数据是否正确
	if loadedConfig.Dest["save-test"].Ssh != "save@example.com" {
		t.Errorf("expected SSH to be 'save@example.com', got '%s'", loadedConfig.Dest["save-test"].Ssh)
	}
	if loadedConfig.Dest["save-test"].Region != "ap-shanghai" {
		t.Errorf("expected region to be 'ap-shanghai', got '%s'", loadedConfig.Dest["save-test"].Region)
	}
	if loadedConfig.Dest["save-test"].InstanceId != "ins-save456" {
		t.Errorf("expected instance ID to be 'ins-save456', got '%s'", loadedConfig.Dest["save-test"].InstanceId)
	}
}

func TestGetConfigFilePath(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	defer os.Setenv("HOME", originalHomeDir)
	defer os.Setenv("USERPROFILE", originalUserProfile)

	expectedPath := filepath.Join(tempDir, CONFIG_DIR, CONFIG_FILE)

	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if path != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, path)
	}

	// 验证目录和文件是否被创建
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		t.Errorf("config directory was not created: %s", filepath.Dir(path))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("config file was not created: %s", path)
	}
}