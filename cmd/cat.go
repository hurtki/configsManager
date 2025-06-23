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

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:   "cat [key]",
	Short: "Print the content of the config file for a given key",
	Long: `The 'cat' command fetches and displays the full content of the configuration file 
associated with the specified key from the user's saved configs list.

Usage example:
  cm cat myconfig

This command is useful when you want to quickly inspect the contents of a configuration 
file without opening it in an editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: missing required argument 'key'")
			os.Exit(1)
		}
		key := args[0]

		data, err := service.GetFileDataByConfigKey(key)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(data)
	},
}


func init() {
	rootCmd.AddCommand(catCmd)
}
