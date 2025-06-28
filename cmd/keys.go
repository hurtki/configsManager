/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	"os"

	service "github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List all configuration keys",
	Long: `The 'keys' command outputs all available configuration keys stored 
in the user's configuration list.

Usage example:
  cm keys

This command is useful to quickly see what configuration entries exist 
and can be accessed or managed with other commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := service.GetAllKeys()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for i := 0; i < len(keys); i++ {
			fmt.Println(keys[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
