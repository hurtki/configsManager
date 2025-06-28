package store

import (
	"errors"
	"os"
)

// GetFileText gets path and returns file text 
func GetFileText(filepath string) (string, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return "", errors.New("file not found")
	}
	return string(text), nil
}

// FileExists() checks if file exist on the given path
func FileExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && !info.IsDir()
}