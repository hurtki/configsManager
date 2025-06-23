package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Editor string `json:"editor"`
}

const (
	configRelPath = "/.config/configManager.json"
)

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
		file, err = os.Create(configPath)
		if err != nil {
			return AppConfig{}, err
		}

		//writing default config to file
		defaultConfig := AppConfig {
			Editor: "vim",
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(defaultConfig)
		// to the start of the file 
		file.Seek(0, 0)

		if err != nil {
			panic(err)
		}
	}
	
	// defering file to close
	defer file.Close()

	var config AppConfig
	// decoding of config file
	err = json.NewDecoder(file).Decode(&config)
	
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}




