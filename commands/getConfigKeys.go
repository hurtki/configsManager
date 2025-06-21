package commands

import "github.com/hurtki/configManager/utils"

// GetConfigKeys returns all keys in configs list file
func GetConfigKeys() ([]string, error) {
	configs, err := utils.GetConfigsListFile()
	if err != nil {
		return make([]string, 0), err
	}
	keys := make([]string, 0, len(configs)) 
	for k := range configs {
    	keys = append(keys, k)
	}
	return keys, nil
}