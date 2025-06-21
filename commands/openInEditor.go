package commands

import "github.com/hurtki/configManager/utils"

func OpenInEditor(key string) error {
	path, err := GetPathConfigPath(key)
	if err != nil {
		return err
	}
	err = utils.OpenInEditor(path)
	return err
} 