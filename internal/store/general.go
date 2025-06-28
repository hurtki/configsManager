package store

import (
	"errors"
	"os"
)

func GetFileText(filepath string) (string, error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return "", errors.New("file not found")
	}
	return string(text), nil
}

func FileExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && !info.IsDir()
}