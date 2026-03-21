package sync_cmd

import (
	"github.com/hurtki/configsManager/internal/domain"
	sync_services "github.com/hurtki/configsManager/internal/services/sync"
)

type SyncService interface {
	// Authorization
	Auth(provider string) error
	Logout(provider string) error // blank provider param => logout for everyone

	// Pulling
	PullAll() ([]sync_services.SyncResult, error)
	PullOne(key string) sync_services.SyncResult

	// Pushing
	Push(configs []*domain.ConfigObj, force bool) ([]*sync_services.SyncResult, error)
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

type ConfigsListService interface {
	Load() (*domain.ConfigsList, error)
	Save(*domain.ConfigsList) error
	GenerateUniqueKeyForPath(path string) (string, error)
}
