package commands

import (
	u "github.com/hurtki/configManager/utils"
	"errors")

// GetPathConfigPath realizes getting config path from configs list 
// returns error if not found 
func GetPathConfigPath(key string) (string, error) {
	configsList, err := u.GetConfigListFile()
	if err != nil {
		panic(err)
	}
	val, ok := configsList[key]
	if ok {
		return val, nil
	} else {
		return "", errors.New("the config key not found in configs list")
	}	
}

