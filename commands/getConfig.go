package commands

import (
	"errors"

	"github.com/hurtki/configManager/utils"
)

// GetFileDataByConfigKey finding config in configs list and returns file data from its path
// returns error if file or config wasn't found
func GetFileDataByConfigKey(key string) (string, error) {

	value, err := GetPathConfigPath(key)
	if err != nil {
		return "", errors.New("config was not found")
	}

	data, err := utils.GetFileText(value)
	if err != nil {
		return "", errors.New("file from config list not found")
	}
	return data, nil
}