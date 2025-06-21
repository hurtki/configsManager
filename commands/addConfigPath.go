package commands

import (
	"encoding/json"
	"errors"
	"os"

	u "github.com/hurtki/configManager/utils"
)

// addConfigPath creates a value for key in configs list
func AddConfigPath(key string, value string) error {
	data, err := u.GetConfigsListFile()
	if err != nil {
		return err
	}
	data[key] = value
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.New("unable to encode map to json")
	}
	err = os.WriteFile(u.GetAbsoluteConfigsListPath(), jsonData, 0644)

	if err != nil {
		return errors.New("unable to write json to file")
	}
	return err
}