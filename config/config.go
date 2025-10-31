// Package config handles configuration loading and saving for the lucky-go application.
// It provides functionality for loading destination instances and saving configuration data.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const CONFIG_DIR = ".lucky-go"
const CONFIG_FILE = "config.yaml"

// Config represents the main configuration structure for the application.
type Config struct {
	// Dest maps destination names to destination instances
	Dest map[string]DestinationInstance `yaml:"dest"`
}

// DestinationInstance represents a cloud instance with SSH connection details.
type DestinationInstance struct {
	// Ssh contains the SSH connection string
	Ssh string `yaml:"ssh"`
	// Region specifies the cloud region
	Region string `yaml:"region"`
	// InstanceId is the unique identifier for the instance
	InstanceId string `yaml:"instance-id"`
}

// SaveConfig saves the configuration to the config file.
// It marshals the config struct to YAML format and writes it to the config file.
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

// LoadDestinationInstance loads a destination instance from the configuration by name.
// It returns the destination instance and any error encountered while loading the config.
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

// loadConfig loads the configuration from the config file.
// It reads the YAML file and unmarshals it into a Config struct.
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

// getConfigFilePath returns the path to the config file, creating directories if needed.
// It ensures the config directory and file exist, creating them if they don't.
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
