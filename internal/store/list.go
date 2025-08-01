package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	ConfigListRelPath = "/.config/configs_list.json"
	CMConfigRelPath   = "/.config/configManager.json"
)

var DefaultConfigsListMap = map[string]string{
	"cm_config": GetAbsoluteCMConfigPath(),
}

// GetAbsoluteConfigsListPath returns absolute path to configs list
// panics if cannot found home dir using os.UserHomeDir()
func GetAbsoluteConfigsListPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ConfigListRelPath)
}

// GetAbsoluteCMConfigPath returns absolute path to CM config
// panics if cannot found home dir using os.UserHomeDir()
func GetAbsoluteCMConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, CMConfigRelPath)
}

// GetConfig loads file with configs list
// If the file does not exist, it is created blank
func LoadUserConfigs() (map[string]string, error) {
	configPath := GetAbsoluteConfigsListPath()

	// trying to open config
	// if no file creating it
	file, err := os.Open(configPath)
	if err != nil {
		file, err = os.Create(configPath)
		if err != nil {
			return make(map[string]string), err
		}

		jsonData, err := json.MarshalIndent(DefaultConfigsListMap, "", "  ")
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

// WriteConfigsList writes data to configs list
// returns error if unable to write
func WriteConfigsList(data map[string]string) error {
	configPath := GetAbsoluteConfigsListPath()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.New("unable to encode map to json")
	}
	err = os.WriteFile(configPath, jsonData, 0644)

	if err != nil {
		return errors.New("unable to write json to file")
	}
	return err

}
