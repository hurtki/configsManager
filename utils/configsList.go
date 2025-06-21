package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)


const (
	ConfigListRelPath = "/.config/configs_list.json"
)

// GetAbsoluteConfigsListPath returns absolute path to configs list 
// panics if cannot found home dir using os.UserHomeDir()
func GetAbsoluteConfigsListPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ConfigListRelPath)
}

// GetConfig loads file with configs list
// If the file does not exist, it is created blank 
func GetConfigListFile() (map[string]string, error) {
	configPath := GetAbsoluteConfigsListPath()

	// trying to open config
	// if no file creating it
	file, err := os.Open(configPath)
	if err != nil {
		file, err = os.Create(configPath)
		if err != nil {
			return make(map[string]string), err
		}

		//writing default config to file
		defaultConfig := make(map[string]string)

		jsonData, err := json.MarshalIndent(defaultConfig, "", "  ") 
		if err != nil {
			return make(map[string]string), err
		}
		err = os.WriteFile(configPath, jsonData, 0644)
		
		// to the start of the file 
		file.Seek(0, 0)

		if err != nil {
			panic(err)
		}
	}
	
	// defering file to close
	defer file.Close()

	var result map[string]string

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&result)
	if err != nil {
		return make(map[string]string), err
	}
	return result, nil
}

