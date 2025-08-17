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

func (n NoopCloudManager) SetConfig(key string, checksum [32]byte) error {
	return ErrNotAuthenticated
}

func (n NoopCloudManager) RemoveConfig(key string) error {
	return ErrNotAuthenticated
}

func (n NoopCloudManager) SetChecksum(key string, checksum [32]byte) error {
	return ErrNotAuthenticated
}

func (n NoopCloudManager) UpdateConfig(cfg ConfigObj) error {
	return ErrNotAuthenticated
}

func (n NoopCloudManager) DownloadConfig(key string) (*ConfigObj, error) {
	return nil, ErrNotAuthenticated
}
