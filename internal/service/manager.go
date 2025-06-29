package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	editor "github.com/hurtki/configsManager/internal/editor"
	store "github.com/hurtki/configsManager/internal/store"
)

// GetConfigKeys returns all keys in configs list file
func GetAllKeys() ([]string, error) {
	configs, err := store.LoadUserConfigs()
	if err != nil {
		return make([]string, 0), err
	}
	keys := make([]string, 0, len(configs)) 
	for k := range configs {
    	keys = append(keys, k)
	}
	return keys, nil
}



// addConfigPath creates a value for key in configs list
func AddConfig(key string, value string) error {
	data, err := store.LoadUserConfigs()
	if err != nil {
		return err
	}
	data[key] = value

	err = store.WriteConfigsList(data)

	return err
}


// GetPathConfigPath realizes getting config path from configs list 
// returns error if not found 
func GetPathByKey(key string) (string, error) {
	configsList, err := store.LoadUserConfigs()
	if err != nil {
		return "", err
	}
	val, ok := configsList[key]
	if ok {
		return val, nil
	} else {
		return "", errors.New("the config key not found in configs list")
	}	
}


// GetFileDataByConfigKey finding config in configs list and returns file data from its path
// returns error if file or config wasn't found
func GetFileDataByConfigKey(key string) (string, error) {

	value, err := GetPathByKey(key)
	if err != nil {
		return "", errors.New("config was not found")
	}

	data, err := store.GetFileText(value)
	if err != nil {
		return "", errors.New("file from config list was not found")
	}
	return data, nil
}

// OpenByKey opens config file by key 
func OpenByKey(key string) error {
	path, err := GetPathByKey(key)
	if err != nil {
		return err
	}
	cfg, err := store.GetConfig()
	if err != nil {
		return  err
	}
	return editor.OpenInEditor(*cfg.Editor, path)
}

// LoadAppConfig loads app config from storage
func LoadAppConfig() (store.AppConfig, error) {
	return store.GetConfig()
}

// Generates unique key for filepath
// the key won't maths any existing keys in configs list
func GenerateUniqueKeyForPath(path string) (string, error) {
	existingKeys, err := GetAllKeys()
	if err != nil {
		return "", err
	}
	fileName := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	fileName = strings.Split(fileName, ".")[0]
	tempFileName := fileName


	for i := 1;i<1000;i++ {
		if !contains(existingKeys, tempFileName) {
			return tempFileName, nil
		}
		tempFileName = fmt.Sprintf("%s%d", fileName, i)
	} 
	return "", errors.New("could not found a unique key for the path")
}

// contains() checks slice for containing an object
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// GetSTDIn returns stdin from pipe input if it exists; otherwise returns false
func GetSTDIn() (string, bool) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", false
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", false
	}
	clean := strings.TrimSpace(string(data))
	if clean == "" {
		return "", false
	}
	return clean, true
}

func NeedForAskForKeyOverwrting(key string) (bool, error) {
	keysList, err := GetAllKeys()
	if err != nil {
		return false, err
	}
	config, err := store.GetConfig()
	if err != nil {
		return false, err
	}
	
	if contains(keysList, key) && !*config.ForceOverwrite {
		return true, nil
	}
	return false, nil
}

func AskUserYN() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Only y/n:")
		}
	}
}

func NeedForAskForNotExistingPathSaving(path string) (bool, error) {
	config, err := store.GetConfig()
	if err != nil {
		return false, err
	}

	if !*config.ForceAddPath && !store.FileExists(path) {
		return true, nil
	}
	return false, nil
}