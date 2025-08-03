/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package main

import (
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

	rootCmd.Execute()
}
