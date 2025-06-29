package store

import (
	"errors"
	"os"
	"path/filepath"
)

// GetFileText gets path and returns file text
func GetFileText(filepath string) (string, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return "", errors.New("file not found")
	}
	return string(text), nil
}

// FileExistsAbsPath checks if a file exists and returns its absolute path and a bool
func FileExists(path string) (string, bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() {
		return "", false
	}

	return absPath, true
}