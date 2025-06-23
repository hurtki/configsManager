package editor

import (
	"os"
	"os/exec"
	"errors"
)

// OpenInEditor opens file in editor 
func OpenInEditor(editor_cmd, filename string) error {
	cmd := exec.Command(editor_cmd, filename)
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		command := editor_cmd + " " + filename
		return errors.New("error running editor command: " + command)
	}

	return err
}