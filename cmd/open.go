/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"

	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type OpenCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	OsService          services.OsService
	ConfigsListService services.ConfigsListService
}

func (c *OpenCmd) run(cmd *cobra.Command, args []string) error {
	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	appConfig, err := c.AppConfigService.Load()
	if err != nil {
		return nil
	}
	editor := appConfig.Editor

	key := args[0]
	path, ok := configsList.GetPath(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	fmt.Println(path)
	return c.OsService.OpenInEditor(*editor, path)
}

func NewOpenCmd(AppConfig services.AppConfigService,
	ConfigsListService services.ConfigsListService,
	OsService services.OsService,
) OpenCmd {
	openCmd := OpenCmd{
		AppConfigService:   AppConfig,
		OsService:          OsService,
		ConfigsListService: ConfigsListService,
	}

	cmd := &cobra.Command{
		Use:   "open [filename]",
		Short: "Opens a config by name",
		Long: `The 'open' command launches the default text editor specified in the application configuration 
		to open the configuration file associated with the provided key.

		This is useful for quickly editing or viewing the content of config files without manually locating them.

		Usage example:
		cm open myconfig

		This command will open the file linked to 'myconfig' in the configured editor, 
		you can configure it in app config by running 'cm open cm_config'`,
		RunE: openCmd.run,
	}
	openCmd.Command = cmd

	return openCmd
}
