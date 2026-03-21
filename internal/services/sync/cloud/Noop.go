package cloud

import (
	"github.com/hurtki/configsManager/internal/domain"
	sync_services "github.com/hurtki/configsManager/internal/services/sync"
)

type NoopCloudManager struct {
	Error error
}

func (m NoopCloudManager) GetCloudInfo() (*domain.CloudConfigRegistry, error) {
	return nil, m.Error
}

func (m NoopCloudManager) GetChecksum(key string) (*[32]byte, error) {
	return nil, m.Error
}

func (m NoopCloudManager) GetAllKeys() ([]string, error) {
	return nil, m.Error
}

func (m NoopCloudManager) SaveCloudConfigRegistry(reg domain.CloudConfigRegistry) error {
	return m.Error
}

func (m NoopCloudManager) UpdateConfig(cfg domain.ConfigObj) error {
	return m.Error
}

func (m NoopCloudManager) ConcurrentUpdateConfigs(configs []*domain.ConfigObj) ([]*sync_services.SyncResult, error) {
	return nil, m.Error
}

func (m NoopCloudManager) DownloadConfig(key string) (*domain.ConfigObj, error) {
	return nil, m.Error
}
