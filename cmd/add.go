package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type AddCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	InputService       services.InputService
	ConfigsListService services.ConfigsListService
	OsService          services.OsService
}

func (c *AddCmd) run(cmd *cobra.Command, args []string) error {
	// ========================
	// args parsing

	// variants:
	// 1. no args + stdIN => args[0]=stdIN =>
	// 2. one argument + stdIN => key=argument ; value=stdIN
	// 3. one argument(should be a path) + NO stdIN => key=unique key from argument ; value=argument
	// 4. two arguments => key=first argument ; value=second argument

	pipedData, isPipe := c.InputService.GetPipedInput()
	if len(args) < 1 && !isPipe {
		return fmt.Errorf("not enough args")

	}

	appConfig, err := c.AppConfigService.Load()

	if err != nil {
		return err
	}

	var key string
	var value string

	if len(args) < 2 {
		if isPipe && len(args) == 0 {
			key = strings.TrimSuffix(filepath.Base(pipedData), filepath.Ext(pipedData))
			value = pipedData
		} else if isPipe && len(args) == 1 {
			key = args[0]
			value = pipedData
		} else {
			value = args[0]
			key = strings.TrimSuffix(filepath.Base(args[0]), filepath.Ext(args[0]))
		}
	} else {
		key = args[0]
		value = args[1]
	}

	// ========================
	// key and value validating

	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}

	if configsList.HasKey(key) {
		switch *appConfig.IfKeyExists {
		case "default", "ask":
			key, err = c.resolveKeyConflict(key, *appConfig.IfKeyExists)
			if err != nil {
				return err
			}
		case "o":
		case "n":
			var err error
			key, err = c.ConfigsListService.GenerateUniqueKeyForPath(key)
			if err != nil {
				return err
			}
			fmt.Printf("New generated key: %s\n", key)
		}
	}

	path_exists, err := c.OsService.FileExists(value)
	if err != nil {
		return err
	}

	// getting absolute path if the path actually exists
	if path_exists {
		value, err = c.OsService.GetAbsolutePath(value)
		if err != nil {
			return err
		}
	} else {
		err = c.OsService.MakePathAndFile(value)
		if err != nil {
			return err
		}
		value, err = c.OsService.GetAbsolutePath(value)
		if err != nil {
			return err
		}
		fmt.Printf("file created at path: %s\n", value)
	}

	// ========================
	// config adding and saving
	configsList.SetConfig(key, value)

	err = c.ConfigsListService.Save(configsList)

	if err != nil {
		return err
	}
	return nil
}

func NewAddCmd(AppConfigService services.AppConfigService,
	InputService services.InputService,
	ConfigsListService services.ConfigsListService,
	OsService services.OsService,
) *AddCmd {
	addCmd := AddCmd{
		AppConfigService:   AppConfigService,
		InputService:       InputService,
		ConfigsListService: ConfigsListService,
		OsService:          OsService,
	}

	cmd := &cobra.Command{
		Use:   "add [key] [path]",
		Short: "Add a new configuration key with its associated file path",
		Long: `The 'add' command creates a new entry in the user's configuration list,
	associating the provided key with the specified file path.

	Usage examples:
	cm add myconfig /path/to/config/file

	This command is useful to register new configuration files with a key,
	so you can easily reference and manage them later using other commands
	like 'path' or 'cat'.`,
		RunE: addCmd.run,
	}

	addCmd.Command = cmd

	return &addCmd
}

func (c *AddCmd) resolveKeyConflict(key string, mode string) (string, error) {
	basePrompt := fmt.Sprintf(
		"%s already exists in keys, what do you want to do:\n"+
			"o - overwrite\n"+
			"n - create with key %s[num]\n"+
			"q - quit",
		key, key,
	)

	if mode == "default" {
		basePrompt += " (if you don't want to get this message set IfKeyExist to skip it)"
	}
	basePrompt += "\n"

	choice, err := c.InputService.AskUser(basePrompt, []string{"o", "n", "q"})
	if err != nil {
		return "", err
	}

	switch choice {
	case "o":
		return key, nil
	case "n":
		newKey, err := c.ConfigsListService.GenerateUniqueKeyForPath(key)
		if err != nil {
			return "", err
		}
		fmt.Printf("New generated key: %s\n", newKey)
		return newKey, nil
	case "q":
		return "", ErrUserAborted
	default:
		return "", fmt.Errorf("unexpected choice: %s", choice)
	}
}
