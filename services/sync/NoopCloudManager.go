package sync_services

type NoopCloudManager struct {
	Error error
}

func (m NoopCloudManager) GetCloudInfo() (*CloudConfigRegistry, error) {
	return nil, m.Error
}

func (m NoopCloudManager) GetChecksum(key string) (*[32]byte, error) {
	return nil, m.Error
}

func (m NoopCloudManager) GetAllKeys() ([]string, error) {
	return nil, m.Error
}

func (m NoopCloudManager) SaveCloudConfigRegistry(reg CloudConfigRegistry) error {
	return m.Error
}

func (m NoopCloudManager) UpdateConfig(cfg ConfigObj) error {
	return m.Error
}

func (m NoopCloudManager) ConcurrentUpdateConfigs(configs []*ConfigObj) ([]*SyncResult, error) {
	return nil, m.Error
}

func (m NoopCloudManager) DownloadConfig(key string) (*ConfigObj, error) {
	return nil, m.Error
}
