// Package config 处理 lucky-go 应用程序的配置加载和保存。
// 它提供加载目标实例和保存配置数据的功能。
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const CONFIG_DIR = ".lucky-go"
const CONFIG_FILE = "config.yaml"

// Config 表示应用程序的主要配置结构。
type Config struct {
	// Dest 将目标名称映射到目标实例
	Dest map[string]DestinationInstance `yaml:"dest"`
}

// DestinationInstance 表示具有SSH连接详细信息的云实例。
type DestinationInstance struct {
	// Ssh 包含SSH连接字符串
	Ssh string `yaml:"ssh"`
	// Region 指定云区域
	Region string `yaml:"region"`
	// InstanceId 是实例的唯一标识符
	InstanceId string `yaml:"instance-id"`
}

// SaveConfig 将配置保存到配置文件中。
// 它将配置结构体编组为YAML格式并写入配置文件。
func (config Config) SaveConfig() error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

// LoadDestinationInstance 按名称从配置中加载目标实例。
// 它返回目标实例和加载配置时遇到的任何错误。
func LoadDestinationInstance(dest string) (*DestinationInstance, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	res, ok := config.Dest[dest]
	if !ok {
		return nil, fmt.Errorf("配置中不存在目标 %v", dest)
	}

	return &res, nil
}

// loadConfig 从配置文件中加载配置。
// 它读取YAML文件并将其解组到Config结构体中。
func loadConfig() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// getConfigFilePath 返回配置文件的路径，必要时创建目录。
// 它确保配置目录和文件存在，如果不存在则创建它们。
func getConfigFilePath() (string, error) {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, CONFIG_DIR)

	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(configDir, CONFIG_FILE)

	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		f, err := os.Create(configFile)

		if err != nil {
			return "", err
		}
		defer f.Close()
	}

	return configFile, nil
}
