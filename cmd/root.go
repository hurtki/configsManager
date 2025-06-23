package cmd

import (
	"fmt"
	"os"

	store "github.com/hurtki/configManager/internal/store"
	service "github.com/hurtki/configManager/internal/service"
	editor "github.com/hurtki/configManager/internal/editor"
)

func Execute() {
	// cm config manager 
	cfg, err := store.GetConfig()
	if err != nil {
		fmt.Println("error", err.Error())
		panic(err)
	}
	

	// getting command arguments
	args := os.Args
	argsLen := len(args)
	first_arg := ""
	if argsLen < 2{
		first_arg = ""
	} else {
		first_arg = args[1]
	}
	
	// commads handling
	switch first_arg{

		case "--path", "-p":
			if argsLen < 3 {
				fmt.Println("not enough args")
				os.Exit(1)
			}
			key := args[2]
			// service.Get
			value, err := service.GetPathByKey(key)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println(value)
		
		case "--cat", "-c":
			if argsLen < 3 {
				fmt.Println("not enough args")
				os.Exit(1)
			}
			key := args[2]

			data, err := service.GetFileDataByConfigKey(key)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(data)
		
		case "--add", "-a":
			if argsLen < 4 {
				fmt.Println("not enough args")
				os.Exit(1)
			}
			key := args[2]
			value := args[3]

			err := service.AddConfig(key, value)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("successfully added new config")

		case "-n", "--names", "--keys":
			keys, err := service.GetAllKeys()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			for i := 0; i < len(keys); i++ {
				fmt.Println(keys[i])
			}

		default:
			if len(first_arg) == 0{
				fmt.Println("Config manager: for more 'cm -h'")
				return
			}
			err := editor.OpenInEditor(cfg.Editor, first_arg)
			if err != nil {
				fmt.Println("error", err)
			}

	}
}