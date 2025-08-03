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

	if len(args) < 1 {
		data, ok := c.InputService.GetPipedInput()
		if ok {
			args = append(args, data)
		} else {
			return fmt.Errorf("not enough args")
		}
	}

	appConfig, err := c.AppConfigService.Load()

	if err != nil {
		return err
	}

	var key string
	var value string

	if len(args) < 2 {
		data, ok := c.InputService.GetPipedInput()
		if ok {
			value = data
			key = args[0]
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
		case "default":
			prompt := fmt.Sprintf(
				"%s already exists in keys, what do you want to do:\n"+
					"o - overwrite\n"+
					"n - create with key %s[num]\n"+
					"q - quit (if you don't want to get this message set IfKeyExist to skip it)\n",
				key, key,
			)

			choice, err := c.InputService.AskUser(prompt, []string{"o", "n", "q"})

			if err != nil {
				return err
			}
			switch choice {
			case "o":
			case "n":
				key, err = c.ConfigsListService.GenerateUniqueKeyForPath(value)
				if err != nil {
					return err
				}
				fmt.Printf("New generated key: %s", key)
			case "q":
				return nil
			}
		case "o":
		case "n":
			key, err = c.ConfigsListService.GenerateUniqueKeyForPath(value)
			if err != nil {
				return err
			}
		case "ask":
			prompt := fmt.Sprintf(
				"%s already exists in keys, what do you want to do:\n"+
					"o - overwrite\n"+
					"n - create with key %s[num]\n"+
					"q - quit \n",
				key, key,
			)
			choice, err := c.InputService.AskUser(prompt, []string{"o", "n", "q"})

			if err != nil {
				return err
			}
			switch choice {
			case "o":
			case "n":
				key, err = c.ConfigsListService.GenerateUniqueKeyForPath(value)
				if err != nil {
					return err
				}
				fmt.Printf("New generated key: %s", key)
			case "q":
				return nil
			}
		}
	}

	path_exists, err := c.OsService.FileExists(value)
	if err != nil {
		return err
	}

	if !*appConfig.ForceAddPath && !path_exists {
		result, err := c.InputService.AskUser("The path you want to assign is not real, want to continue?",
			[]string{"y", "n"},
		)
		if err != nil {
			return err
		}
		if result == "n" {
			return nil
		}
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
