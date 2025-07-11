/*
Copyright © 2025 Alexey asboba2101@gmail.com >

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cm",
	Short: "A CLI tool to manage configuration file paths by keys",
	Long: `ConfigManager is a simple and efficient CLI application that helps you 
manage your configuration files by associating keys with file paths.

With configManager, you can:
- Add new config entries with keys and paths
- Retrieve config paths by keys
- View content of config files
- List all stored config keys
- Open config files in your preferred editor

Example usage:
  cm add myconfig /path/to/config
  cm path myconfig
  cm cat myconfig

Built with Cobra CLI library to provide a powerful and user-friendly command line interface.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.configManager.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


