// Entry point of the configManager CLI.
// All command-line logic is handled in the `cmd` package.
//
// Author: hurtki
// https://github.com/hurtki/configManager
package main

import "github.com/hurtki/configManager/cmd"

func main() {
	cmd.Execute()
}
