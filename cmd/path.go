/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"
	"os"

	service "github.com/hurtki/configManager/internal/service"
	"github.com/spf13/cobra"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path [key]",
	Short: "Retrieve the file path associated with a configuration key",
	Long: `The 'path' command fetches the absolute file path stored under the given configuration key 
from the user's saved configuration list.

Usage examples:
  cm path myconfig      # Prints the file path associated with "myconfig"
  cm path cm_config     # Prints the default configManager config path

This command helps you quickly locate configuration files by their keys, making management easier 
and more efficient.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("not enough args: please specify the config key")
			os.Exit(1)
		}
		key := args[0]
		value, err := service.GetPathByKey(key)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(value)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
