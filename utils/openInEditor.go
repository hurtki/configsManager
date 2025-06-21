package utils

import (
	"os"
	"os/exec"

	"github.com/hurtki/configManager/config"
)

func OpenInEditor(filename string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	
	cmd := exec.Command(cfg.Editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}