/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hurtki/configsManager/internal/service"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [key]",
	Short: "Delete a specific key from the configuration list",
	Long: `The 'rm' command allows you to delete a configuration key from your local config file.
You must specify the key name you want to remove.

Example:
  cm rm editor

This command will remove the 'editor' key and its value from your configuration file.
If the specified key does not exist, an error will be shown and no changes will be made.

Use this command with caution — deleted keys cannot be recovered unless you re-add them manually.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("not enough args")
			os.Exit(1)
		}
		key := args[0]
		err := service.RemoveConfig(key)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
