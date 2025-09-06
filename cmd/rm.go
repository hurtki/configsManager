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
	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	for _, key := range args {
		configsList.RemoveConfig(key)
	}
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
		Use:   "rm key1 key2 ...",
		Short: "Delete keys from the configuration list",
		Long: `The 'rm' command allows you to delete a cuple or only one key from your local config file.
	You must specify the key/s name/s you want to remove.

	Example:
	cm rm nginx django

	This command will remove the 'nginx' and 'django' keys and its values from your configuration file.
	If the some of the keys does not exist, an error won't be shown and other keys will be deleted correctly.
	
	Use this command with caution — deleted keys cannot be recovered unless you re-add them manually.`,
		RunE: rmCmd.run,
	}
	rmCmd.Command = cmd

	return rmCmd
}
