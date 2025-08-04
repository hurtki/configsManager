/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package main

import (
	"fmt"
	"os"

	"github.com/hurtki/configsManager/cmd"
	"github.com/hurtki/configsManager/services"
)

func main() {
	// making dependencies
	AppConfigService := services.NewAppConfigServiceImpl()
	StdInputService := services.NewStdInputService()
	ConfigsListService := services.NewConfigsListServiceImpl()
	OsService := services.NewOsServiceImpl()

	rootCmd := cmd.NewRootCmd(AppConfigService, StdInputService, ConfigsListService, OsService)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
