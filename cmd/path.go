/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type PathCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	ConfigsListService services.ConfigsListService
}

func (c *PathCmd) run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough args: please specify the config key")
	}
	key := args[0]
	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	path, ok := configsList.GetPath(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	fmt.Println(path)
	return nil
}

func NewPathCmd(AppConfig services.AppConfigService, ConfigsListService services.ConfigsListService) PathCmd {
	pathCmd := PathCmd{
		AppConfigService:   AppConfig,
		ConfigsListService: ConfigsListService,
	}

	cmd := &cobra.Command{
		Use:   "path [key]",
		Short: "Retrieve the file path associated with a configuration key",
		Long: `The 'path' command fetches the absolute file path stored under the given configuration key 
from the user's saved configuration list.

Usage examples:
  cm path myconfig      # Prints the file path associated with "myconfig"
  cm path cm_config     # Prints the default configManager config path

This command helps you quickly locate configuration files by their keys, making management easier 
and more efficient.`,
		RunE: pathCmd.run,
	}
	pathCmd.Command = cmd

	return pathCmd
}
