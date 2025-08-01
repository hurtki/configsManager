/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	service "github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
)

type OpenCmd struct {
	Command *cobra.Command
	AppConfig *service.AppConfig
}

func (k *OpenCmd) run(cmd *cobra.Command, args []string) error {
	err := service.OpenByKey(args[0])
	if err != nil {
		return err
	}
	return nil
}


func NewOpenCmd(AppConfig *service.AppConfig) OpenCmd {
	openCmd := OpenCmd{
		AppConfig: AppConfig,
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
