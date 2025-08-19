package sync_services

import (
	"io"
	"path/filepath"
)

type ConfigObj struct {
	KeyName        string
	FileName       string // with extension
	Content        []byte
	DeterminedPath string
}

func NewConfigObj(file io.Reader, path, key string) (*ConfigObj, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &ConfigObj{
		KeyName:        key,
		FileName:       filepath.Base(path),
		Content:        content,
		DeterminedPath: filepath.Dir(path),
	}, nil
}
