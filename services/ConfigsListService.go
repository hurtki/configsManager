package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigsListService interface {
	Load() (*ConfigsList, error)
	Save(*ConfigsList) error
	GenerateUniqueKeyForPath(path string) (string, error)
}

const (
	configsListFileName = "configs_list.json"
)

type ConfigsListServiceImpl struct {
}

func NewConfigsListServiceImpl() *ConfigsListServiceImpl {
	return &ConfigsListServiceImpl{}
}

func (c *ConfigsListServiceImpl) getConfigsListDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, AppDir), nil
}

func (c *ConfigsListServiceImpl) getConfigsListPath() (string, error) {
	configDir, err := c.getConfigsListDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configsListFileName), nil
}

func (s *ConfigsListServiceImpl) Load() (*ConfigsList, error) {
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
		return GetDefaultConfigsList(AppConfigPath), nil
	}

	file, err := os.ReadFile(configsListPath)

	if err != nil {
		return nil, err
	}

	var configsList ConfigsList

	if err := json.Unmarshal(file, &configsList.configs); err != nil {
		return nil, err
	}

	return &configsList, nil
}

func (s *ConfigsListServiceImpl) Save(cfgList *ConfigsList) error {
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

	jsonData, err := json.MarshalIndent(cfgList.configs, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configsListPath, jsonData, 0644)

	return err
}

func (s *ConfigsListServiceImpl) GenerateUniqueKeyForPath(path string) (string, error) {
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
