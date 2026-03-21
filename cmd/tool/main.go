package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hurtki/configsManager/internal/cmd"
	services "github.com/hurtki/configsManager/internal/services"
	sync_services "github.com/hurtki/configsManager/internal/services/sync"
	"github.com/hurtki/configsManager/internal/services/sync/auth"
	"github.com/hurtki/configsManager/internal/services/sync/cloud"
)

func main() {
	// making dependencies
	appConfigService := services.NewAppConfigService()
	stdInputService := services.NewStdInputService()
	configsListService := services.NewConfigsListService()
	osService := services.NewOsService()

	tokenStore := auth.NewTokenStoreImpl()
	authManager := auth.NewAuthManager(tokenStore)

	var cloudManager sync_services.CloudManager
	token, err := authManager.GetToken("dropbox")

	if err == nil {
		cloudManager = cloud.NewCloudManager(token)
	} else {
		cloudManager = cloud.NoopCloudManager{Error: err}
	}

	syncService := sync_services.NewSyncService(authManager, cloudManager)

	rootCmd := cmd.NewRootCmd(appConfigService,
		stdInputService,
		configsListService,
		osService,
		syncService,
	)

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, cmd.ErrUserAborted) {
			os.Exit(3)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Println("For help with this command: 'cm help'")
		os.Exit(1)
	}
}
