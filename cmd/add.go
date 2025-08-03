/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type AddCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	InputService       services.InputService
	ConfigsListService services.ConfigsListService
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
			key, err = c.ConfigsListService.GenerateUniqueKeyForPath(value)
			if err != nil {
				return err
			}
			fmt.Println("key assighned to your path:", key)
		}
	} else {
		key = args[0]
		value = args[1]
	}

	// ========================
	// key and value validating

	if *appConfig.ForceOverwrite {
		accept, err := c.InputService.AskUserYN("The key you want to assign already exist, want to overwrite?")
		if err != nil {
			return err
		}
		if !accept {
			return fmt.Errorf("operation cancelled by user")
		}
	}

	if *appConfig.ForceAddPath {
		accept, err := c.InputService.AskUserYN("The path you want to assign is not real, want to continue?")
		if err != nil {
			return err
		}
		if !accept {
			return fmt.Errorf("operation cancelled by user")
		}
	}

	// ========================
	// config adding

	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
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
) *AddCmd {
	addCmd := AddCmd{
		AppConfigService:   AppConfigService,
		InputService:       InputService,
		ConfigsListService: ConfigsListService,
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
