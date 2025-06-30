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
		// ========================
		// args parsing

		// variants:
		// 1. no args + stdIN => args[0]=stdIN => 
		// 2. one argument + stdIN => key=argument ; value=stdIN
		// 3. one argument(should be a path) + NO stdIN => key=unique key from argument ; value=argument
		// 4. two arguments => key=first argument ; value=second argument

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
		
		// ========================
		// key and value validating 

		shouldAskOverwrite, err := service.ShouldConfirmOverwrite(key)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if shouldAskOverwrite {
			fmt.Println("The key you want to assign already exist, want to overwrite? y/n")
			accept := service.AskUserYN(os.Stdin)
			if !accept {
				os.Exit(1)
			}
		}

		shouldAskPathConfirmation, err := service.ShouldConfirmInvalidPath(value)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		
		if shouldAskPathConfirmation {
			fmt.Println("The path you want to assign is not real, want to continue? y/n")
			accept := service.AskUserYN(os.Stdin)
			if !accept {
				os.Exit(1)
			}
		}
		
		// ========================
		// config adding
		
		err = service.AddConfig(key, value)
		
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
