package services

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hurtki/configsManager/internal/config"
)

type AppConfigService struct {
}

func NewAppConfigService() *AppConfigService {
	return &AppConfigService{}
}

const (
	AppDir         = ".config/configsManager"
	configFileName = "configsManager.json"
)

func (c *AppConfigService) getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppDir), nil
}

func (c *AppConfigService) getConfigPath() (string, error) {
	configDir, err := c.getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

func (c *AppConfigService) Load() (*config.AppConfig, error) {
	configPath, err := c.getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config.NewDefaultAppConfig(), nil
	}

	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var config config.AppConfig

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// validating and fil
	if config.ValidateAppConfig() {
		if err := c.Save(&config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func (c *AppConfigService) Save(cfg *config.AppConfig) error {
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
