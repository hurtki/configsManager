package sync_services

type NoopCloudManager struct{}

func (n NoopCloudManager) GetCloudInfo() (*CloudConfigRegistry, error) {
	return nil, ErrNotAuthenticated
}

func (n NoopCloudManager) GetChecksum(key string) (*[32]byte, error) {
	return nil, ErrNotAuthenticated
}

func (n NoopCloudManager) GetAllKeys() ([]string, error) {
	return nil, ErrNotAuthenticated
}

func (n NoopCloudManager) SaveCloudConfigRegistry(reg CloudConfigRegistry) error {
	return ErrNotAuthenticated
}

func (n NoopCloudManager) UpdateConfig(cfg ConfigObj) error {
	return ErrNotAuthenticated
}

func (m NoopCloudManager) ConcurrentUpdateConfigs(configs []*ConfigObj) ([]*SyncResult, error) {
	return nil, ErrNotAuthenticated
}

func (n NoopCloudManager) DownloadConfig(key string) (*ConfigObj, error) {
	return nil, ErrNotAuthenticated
}
