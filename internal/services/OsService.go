package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type OsService struct{}

func NewOsService() *OsService {
	return &OsService{}
}

// MakePathAndFile creates all directories in the given path if they do not exist,
// and then creates or truncates the file at the end of the path.
// Returns an error if any operation fails.
func (s *OsService) MakePathAndFile(path string) error {
	dir := filepath.Dir(path)
	// creating all the folders
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// create or
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("error closing file")
		}
	}()
	return nil
}

func (s *OsService) GetFileData(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *OsService) OpenInEditor(editor, path string) error {
	cmd := exec.Command(editor, path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func (s *OsService) FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *OsService) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}
func (s *OsService) WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (s *OsService) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}
