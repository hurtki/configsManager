package cmd

import (
	"fmt"
	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type RmCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	ConfigsListService services.ConfigsListService
}

func (c *RmCmd) run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough args")
	}
	key := args[0]
	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	configsList.RemoveConfig(key)
	if err := c.ConfigsListService.Save(configsList); err != nil {
		return err
	}
	return nil
}

func NewRmCmd(AppConfig services.AppConfigService, ConfigsListService services.ConfigsListService) RmCmd {
	rmCmd := RmCmd{
		AppConfigService:   AppConfig,
		ConfigsListService: ConfigsListService,
	}

	cmd := &cobra.Command{
		Use:   "rm [key]",
		Short: "Delete a specific key from the configuration list",
		Long: `The 'rm' command allows you to delete a configuration key from your local config file.
	You must specify the key name you want to remove.

	Example:
	cm rm editor

	This command will remove the 'editor' key and its value from your configuration file.
	If the specified key does not exist, an error will be shown and no changes will be made.

	Use this command with caution — deleted keys cannot be recovered unless you re-add them manually.`,
		RunE: rmCmd.run,
	}
	rmCmd.Command = cmd

	return rmCmd
}
