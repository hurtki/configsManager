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