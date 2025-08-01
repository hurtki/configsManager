/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	service "github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
)

type CatCmd struct {
	Command   *cobra.Command
	AppConfig *service.AppConfig
}

// catCmd represents the cat command
func (c *CatCmd) run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Printf("missing required argument 'key'")
	}
	key := args[0]

	data, err := service.GetFileDataByConfigKey(key)
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}

func NewCatCmd(AppConfig *service.AppConfig) *CatCmd {
	catCmd := CatCmd{
		AppConfig: AppConfig,
	}

	cmd := &cobra.Command{
		Use:   "cat [key]",
		Short: "Print the content of the config file for a given key",
		Long: `The 'cat' command fetches and displays the full content of the configuration file 
	associated with the specified key from the user's saved configs list.

	Usage example:
	cm cat myconfig

	This command is useful when you want to quickly inspect the contents of a configuration 
	file without opening it in an editor.`,
		RunE: catCmd.run,
	}

	catCmd.Command = cmd

	return &catCmd
}
