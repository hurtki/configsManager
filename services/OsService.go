package services

import (
	"os"
	"os/exec"
	"path/filepath"
)

type OsService interface {
	GetFileData(path string) ([]byte, error)
	OpenInEditor(editor, path string) error
	FileExists(path string) (bool, error)
	GetAbsolutePath(path string) (string, error)
}

type OsServiceImpl struct{}

func NewOsServiceImpl() *OsServiceImpl {
	return &OsServiceImpl{}
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