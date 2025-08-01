/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	service "github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
)

type KeysCmd struct {
	Command   *cobra.Command
	AppConfig *service.AppConfig
}

func (k *KeysCmd) run(cmd *cobra.Command, args []string) error {
	keys, err := service.GetAllKeys()
	if err != nil {
		return err
	}
	for i := 0; i < len(keys); i++ {
		fmt.Println(keys[i])
	}
	return nil
}

func NewKeysCmd(AppConfig *service.AppConfig) KeysCmd {
	keysCmd := KeysCmd{
		AppConfig: AppConfig,
	}

	cmd := &cobra.Command{
		Use:   "keys",
		Short: "List all configuration keys",
		Long: `The 'keys' command outputs all available configuration keys stored 
		in the user's configuration list.

		Usage example:
		cm keys

		This command is useful to quickly see what configuration entries exist 
		and can be accessed or managed with other commands.`,
		RunE: keysCmd.run,
	}
	keysCmd.Command = cmd

	return keysCmd
}
