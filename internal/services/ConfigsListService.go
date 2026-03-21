package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hurtki/configsManager/internal/domain"
)

const (
	configsListFileName = "configs_list.json"
)

type ConfigsListService struct {
}

func NewConfigsListService() *ConfigsListService {
	return &ConfigsListService{}
}

func (c *ConfigsListService) getConfigsListDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppDir), nil
}

func (c *ConfigsListService) getConfigsListPath() (string, error) {
	configDir, err := c.getConfigsListDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configsListFileName), nil
}

func (s *ConfigsListService) Load() (*domain.ConfigsList, error) {
	configsListPath, err := s.getConfigsListPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configsListPath); os.IsNotExist(err) {
		configsListDir, err := s.getConfigsListDir()
		if err != nil {
			return nil, err
		}
		AppConfigPath := filepath.Join(configsListDir, configFileName)
		return domain.GetDefaultConfigsList(AppConfigPath), nil
	}

	file, err := os.ReadFile(configsListPath)

	if err != nil {
		return nil, err
	}

	var cfgs map[string]string

	if err := json.Unmarshal(file, &cfgs); err != nil {
		return nil, err
	}

	return domain.NewConfigsList(cfgs), nil
}

func (s *ConfigsListService) Save(cfgList *domain.ConfigsList) error {
	configsListDir, err := s.getConfigsListDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(configsListDir, 0755); err != nil {
		return err
	}

	configsListPath, err := s.getConfigsListPath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cfgList.Configs, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configsListPath, jsonData, 0644)

	return err
}

func (s *ConfigsListService) GenerateUniqueKeyForPath(path string) (string, error) {
	cfgList, err := s.Load()
	if err != nil {
		return "", err
	}

	baseName := filepath.Base(path)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	tempKey := baseName

	for i := 1; i < 1000; i++ {
		if !cfgList.HasKey(tempKey) {
			return tempKey, nil
		}
		tempKey = fmt.Sprintf("%s%d", baseName, i)
	}

	return "", errors.New("could not find a unique key for the path")
}
