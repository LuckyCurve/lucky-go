package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const CONFIG_DIR = ".lucky-go"
const CONFIG_FILE = "config.yaml"

type (
	Config struct {
		Dest map[string]DestinationInstance `yaml:"dest"`
	}

	DestinationInstance struct {
		Ssh        string `yaml:"ssh"`
		Region     string `yaml:"region"`
		InstanceId string `yaml:"instance-id"`
	}
)

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

func LoadDestinationInstance(dest string) (*DestinationInstance, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	res, ok := config.Dest[dest]
	if !ok {
		return nil, fmt.Errorf("not destination %v in config", dest)
	}

	return &res, nil
}

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
