package services

import (
	"os"
	"os/exec"
)

type OsService interface {
	GetFileData(path string) ([]byte, error)
	OpenInEditor(editor, path string) error
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
