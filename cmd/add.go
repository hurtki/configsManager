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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [key] [path]",
	Short: "Add a new configuration key with its associated file path",
	Long: `The 'add' command creates a new entry in the user's configuration list,
associating the provided key with the specified file path.

Usage examples:
  cm add myconfig /path/to/config/file

This command is useful to register new configuration files with a key,
so you can easily reference and manage them later using other commands
like 'path' or 'cat'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			data, ok := service.GetSTDIn()
			if ok {
				args = append(args, data)
			} else {
				fmt.Println("not enough args")
				os.Exit(1)
			}
		}
		var key string
		var value string
		var err error

		if len(args) < 2 {
			data, ok := service.GetSTDIn()
			if ok {
				value = data
				key = args[0]
			} else {
				value = args[0]
				key, err = service.GenerateUniqueKeyForPath(value)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				fmt.Println("key assighned to your path:", key)
			}
		} else {
			key = args[0]
			value = args[1]
		}
		err = service.AddConfig(key, value)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("successfully added new config")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
