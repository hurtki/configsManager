/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package cmd

import (
	"fmt"

	"github.com/hurtki/configManager/internal/service"
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [filename]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := service.OpenByKey(args[0])
		if err != nil {
			fmt.Println("error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
