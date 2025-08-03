/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	"github.com/hurtki/configsManager/services"
	"github.com/spf13/cobra"
)

type KeysCmd struct {
	Command            *cobra.Command
	AppConfigService   services.AppConfigService
	ConfigsListService services.ConfigsListService
}

func (c *KeysCmd) run(cmd *cobra.Command, args []string) error {

	configsList, err := c.ConfigsListService.Load()
	if err != nil {
		return err
	}
	keys := configsList.GetAllKeys()

	for i := 0; i < len(keys); i++ {
		fmt.Println(keys[i])
	}
	return nil
}

func NewKeysCmd(AppConfig services.AppConfigService,
	ConfigsListService services.ConfigsListService,
) KeysCmd {
	keysCmd := KeysCmd{
		AppConfigService:   AppConfig,
		ConfigsListService: ConfigsListService,
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
