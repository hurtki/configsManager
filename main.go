/*
Copyright © 2025 Alexey asboba2101@gmail.com >
*/
package main

import (
	"github.com/hurtki/configsManager/cmd"
	"github.com/hurtki/configsManager/internal/service"
)

func main() {
	// making dependencies 
	OsAppConfig := service.OsAppConfig {
	}
	
	rootCmd := cmd.NewRootCmd(&OsAppConfig)
	
	rootCmd.Execute()
}
