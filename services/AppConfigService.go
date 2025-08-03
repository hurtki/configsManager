package services

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfigService interface {
	Load() (*AppConfig, error)
	Save(*AppConfig) error
}

type AppConfigServiceImpl struct {
}

func NewAppConfigServiceImpl() *AppConfigServiceImpl {
	return &AppConfigServiceImpl{}
}

const (
	AppDir         = ".config/configsManager"
	configFileName = "configsManager.json"
)

func (c *AppConfigServiceImpl) getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppDir), nil
}

func (c *AppConfigServiceImpl) getConfigPath() (string, error) {
	configDir, err := c.getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

func (c *AppConfigServiceImpl) Load() (*AppConfig, error) {
	configPath, err := c.getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &defaultConfig, nil
	}

	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var config AppConfig

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// validating and fil
	if config.validateAppConfig() {
		if err := c.Save(&config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func (c *AppConfigServiceImpl) Save(cfg *AppConfig) error {
	configDir, err := c.getConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath, err := c.getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
