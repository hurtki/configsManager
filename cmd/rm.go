/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	service "github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
	"fmt"
)

type RmCmd struct {
	Command *cobra.Command
	AppConfig *service.AppConfig
}

func (k *RmCmd) run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough args")
	}
	key := args[0]
	err := service.RemoveConfig(key)
	if err != nil {
		return err
	}
	return nil
}


func NewRmCmd(AppConfig *service.AppConfig) RmCmd {
	rmCmd := RmCmd{
		AppConfig: AppConfig,
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
