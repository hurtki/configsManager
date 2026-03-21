package cmd

import (
	"github.com/hurtki/configsManager/internal/config"
	"github.com/hurtki/configsManager/internal/domain"
)

type AppConfigService interface {
	Load() (*config.AppConfig, error)
	Save(*config.AppConfig) error
}

type ConfigsListService interface {
	Load() (*domain.ConfigsList, error)
	Save(*domain.ConfigsList) error
	GenerateUniqueKeyForPath(path string) (string, error)
}

type OsService interface {
	GetFileData(path string) ([]byte, error)
	OpenInEditor(editor, path string) error
	FileExists(path string) (bool, error)
	GetAbsolutePath(path string) (string, error)
	MakePathAndFile(path string) error
	WriteFile(path string, data []byte) error
	GetHomeDir() (string, error)
}

type InputService interface {
	AskUser(prompt string, options []string) (string, error)
	GetPipedInput() (string, bool)
}
