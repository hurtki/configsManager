package cmd

import (
	"fmt"
	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type InitCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	ConfigsListService services.ConfigsListService
}

func (c *InitCmd) run(cmd *cobra.Command, args []string) error {

	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	appConfig, err := c.AppConfigService.Load()
	if err != nil {
		return err
	}

	err = c.ConfigsListService.Save(configsList)
	if err != nil {
		return err
	}
	err = c.AppConfigService.Save(appConfig)
	if err != nil {
		return err
	}

	fmt.Printf("Initialized cm folder with configs on path %s\n", services.AppDir)

	return nil
}

func NewInitCmd(AppConfig services.AppConfigService,
	ConfigsListService services.ConfigsListService,
) InitCmd {
	InitCmd := InitCmd{
		AppConfigService:   AppConfig,
		ConfigsListService: ConfigsListService,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes anapp folder",
		Long: `The 'init' command initializes the configuration management folder structure 
and creates default configuration files if they do not exist.

This is the first command that should be run before using other commands 
like 'add', 'open', or 'edit', since it sets up the required files and directories.

Usage example:
  cm init

After running this command, your system will be ready to store and manage
configurations via CM tool.`,
		RunE: InitCmd.run,
	}
	InitCmd.Command = cmd

	return InitCmd
}
