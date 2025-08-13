package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type OsService interface {
	GetFileData(path string) ([]byte, error)
	OpenInEditor(editor, path string) error
	FileExists(path string) (bool, error)
	GetAbsolutePath(path string) (string, error)
	MakePathAndFile(path string) error
}

type OsServiceImpl struct{}

func NewOsServiceImpl() *OsServiceImpl {
	return &OsServiceImpl{}
}

// MakePathAndFile creates all directories in the given path if they do not exist,
// and then creates or truncates the file at the end of the path.
// Returns an error if any operation fails.
func (s *OsServiceImpl) MakePathAndFile(path string) error {
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
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %d\n", err)
		}
	}()

	return nil
}

func (s *OsServiceImpl) GetFileData(path string) ([]byte, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func (s *OsServiceImpl) OpenInEditor(editor, path string) error {
	cmd := exec.Command(editor, path)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func (s *OsServiceImpl) FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *OsServiceImpl) GetAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}
