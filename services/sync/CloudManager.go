package sync_services

import (
	"crypto/sha256"
	"encoding/json"
)

type CloudManager interface {
	GetCloudInfo() ([]CloudConfigInfo, error)
	GetConfigInfo(key string) (*CloudConfigInfo, error)
	GetAllKeys() []string

	SaveCloudInfo([]CloudConfigInfo) error
	AddConfigInfo(CloudConfigInfo) error
	RemoveConfigInfo(key string) error
	UpdateChecksum(key string, checksum [32]byte) error

	UploadConfig(ConfigObj) error
	DownloadConfig(key string) (*ConfigObj, error)
}

type CloudManagerImpl struct {
	Provider Provider
}

func (m *CloudManagerImpl) GetCloudInfo() ([]CloudConfigInfo, error) {
	file, err := m.Provider.Download(cloudManagerFileName)
	if err == ErrFileDoesntExists {
		defaultCloudManagerConfig := DefaultCloudManagerConfigFile()
		defaultCloudManagaerFile, err := json.Marshal(defaultCloudManagerConfig)
		if err != nil {
			return nil, err
		}
		if err := m.Provider.Upload(cloudManagerFileName, defaultCloudManagaerFile); err != nil {
			return nil, nil
		}
		return defaultCloudManagerConfig, nil
	} else if err != nil {
		return nil, nil
	}
	var CloudManagerConfig []CloudConfigInfo
	if err := json.Unmarshal(file, &CloudManagerConfig); err != nil {
		return nil, err
	}
	return CloudManagerConfig, nil
}

func (m *CloudManagerImpl) UpdateConfig(cci CloudConfigInfo, co ConfigObj) error {
	CheckSum := sha256.Sum256(co.Content)
	cci.Checksum = CheckSum

	m.Provider.Upload(co.ConfigKeyName, co.Content)

	return m.Provider.Upload(co.ConfigKeyName, co.Content)
}
