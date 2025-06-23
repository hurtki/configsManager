package service

import (
	"errors"
	editor "github.com/hurtki/configManager/internal/editor"
	store "github.com/hurtki/configManager/internal/store"
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
		return "", errors.New("file from config list not found")
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
	return editor.OpenInEditor(cfg.Editor, path)
}

// LoadAppConfig loads app config from storage
func LoadAppConfig() (store.AppConfig, error) {
	return store.GetConfig()
}