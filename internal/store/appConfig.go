package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type AppConfig struct {
	// Editor command that command "open" uses
	Editor *string `json:"editor"`
	// Is cm going to overwrite an existing key in configs list
	// if true cm won't ask you
	// if false cm will ask you "If you want to overwrite"
	ForceOverwrite *bool `json:"overwrite_if_exists"`
	// is cm going to add the path if it doesn't exist
	// if true cm won't ask you
	// if false cm will ask you "If you want to add the non existing path"
	ForceAddPath *bool `json:"force_add_path"`
}

const (
	configRelPath = "/.config/configManager.json"
)

// default pointers for structure
var bFalse = func() *bool { b := false; return &b }()
var bEditor = func() *string { b := "vim"; return &b }()

// DEFAULT CONFIG VALUES
var defaultConfig = AppConfig{
	Editor:         bEditor,
	ForceOverwrite: bFalse,
	ForceAddPath:   bFalse,
}

// validateAppConfig() can validate and insert default values
// if validateAppConfig() changed at least one field it returns true
func (cfg *AppConfig) validateAppConfig() bool {
	changed := false

	if cfg.Editor == nil {
		editor := *defaultConfig.Editor
		cfg.Editor = &editor
		changed = true
	}

	if cfg.ForceOverwrite == nil {
		def := *defaultConfig.ForceOverwrite
		cfg.ForceOverwrite = &def
		changed = true
	}

	if cfg.ForceAddPath == nil {
		def := *defaultConfig.ForceAddPath
		cfg.ForceAddPath = &def
		changed = true
	}

	return changed
}

// GetConfig loads the configuration from config_path.
// If the file does not exist, it is created with default settings.
func GetConfig() (AppConfig, error) {
	// получаем путь к конфигу
	home, err := os.UserHomeDir()
	if err != nil {
		return AppConfig{}, err
	}
	configPath := filepath.Join(home, configRelPath)

	// trying to open config
	// if no file creating it
	file, err := os.Open(configPath)

	if err != nil {
		err := saveConfig(configPath, defaultConfig)
		if err != nil {
			return AppConfig{}, errors.New("error opening app config")
		}
		return defaultConfig, nil
	}

	defer file.Close()

	var config AppConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return AppConfig{}, err
	}

	// Валидируем + дописываем поля
	if config.validateAppConfig() {
		if err := saveConfig(configPath, config); err != nil {
			return AppConfig{}, err
		}
	}

	return config, nil
}

// saveConfig(path string, cfg AppConfig) rewrites a config on given path
func saveConfig(path string, cfg AppConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
