package main

import (
	"fmt"
	"os"

	c "github.com/hurtki/configManager/config"
	"github.com/hurtki/configManager/commands"
	"github.com/hurtki/configManager/utils"
)

func main() {
	// cm config manager 
	_, err := c.GetConfig()
	if err != nil {
		fmt.Println("error", err.Error())
		panic(err)
	}
	//fmt.Println("cm config manager was got succsesfully:", cfg.RepoPath)

	// getting command arguments
	args := os.Args
	argsLen := len(args)
	
	// commads handling
	switch args[1]{
	case "--path", "-p":
		if argsLen < 3 {
			fmt.Println("not enough args")
			os.Exit(1)
		}
		key := args[2]

		value, err := commands.GetPathConfigPath(key)
		if err != nil {
			fmt.Println("config was not found")
			os.Exit(1)
		}
		fmt.Println(value)
	

	case "--cat", "-c":
		if argsLen < 3 {
			fmt.Println("not enough args")
			os.Exit(1)
		}
		key := args[2]

		value, err := commands.GetPathConfigPath(key)
		if err != nil {
			fmt.Println("config was not found")
			os.Exit(1)
		}
	
		data, err := utils.GetFileText(value)
		if err != nil {
			fmt.Println("file from config list not found")
			os.Exit(1)
		}
		fmt.Println(data)
	
	





	
	}
}