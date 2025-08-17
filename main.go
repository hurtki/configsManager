/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hurtki/configsManager/cmd"
	"github.com/hurtki/configsManager/services"
	syncServices "github.com/hurtki/configsManager/services/sync"
)

func main() {
	// making dependencies
	AppConfigService := services.NewAppConfigServiceImpl()
	StdInputService := services.NewStdInputService()
	ConfigsListService := services.NewConfigsListServiceImpl()
	OsService := services.NewOsServiceImpl()
	TokenStore := syncServices.NewTokenStoreImpl()
	AuthManager := syncServices.NewAuthManagerImpl(TokenStore)
	SyncService := syncServices.NewSyncServiceImpl(AuthManager)

	rootCmd := cmd.NewRootCmd(AppConfigService, StdInputService, ConfigsListService, OsService, SyncService)

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, cmd.ErrUserAborted) {
			os.Exit(3)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Println("For help with this command: 'cm help'")
		os.Exit(1)
	}
}
