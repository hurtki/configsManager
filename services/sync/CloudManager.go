package sync_services

import (
	"crypto/sha256"
	"encoding/json"
)

const (
	cloudManagerFileName = "cloud_manger.json"
)

type CloudManager interface {
	GetCloudInfo() (*CloudConfigRegistry, error)
	GetChecksum(key string) (*[32]byte, error)
	GetAllKeys() ([]string, error)

	SaveCloudConfigRegistry(CloudConfigRegistry) error
	SetConfig(key string, checksum [32]byte) error
	RemoveConfig(key string) error
	SetChecksum(key string, checksum [32]byte) error

	UpdateConfig(ConfigObj) error
	DownloadConfig(key string) (*ConfigObj, error)
}

type CloudManagerImpl struct {
	Provider Provider
}

func NewCloudManagerImpl(token string) *CloudManagerImpl {
	return &CloudManagerImpl{
		Provider: NewDropboxProvider(token),
	}
}
func (m *CloudManagerImpl) GetCloudInfo() (*CloudConfigRegistry, error) {
	data, err := m.Provider.Download(cloudManagerFileName)
	if err == ErrFileDoesntExist {
		defaultRegistry := CloudConfigRegistry{}
		bytes, err := json.Marshal(defaultRegistry)
		if err != nil {
			return nil, err
		}
		if err := m.Provider.Upload(cloudManagerFileName, bytes); err != nil {
			return nil, err
		}
		return &defaultRegistry, nil
	} else if err != nil {
		return nil, err
	}

	var configRegistry CloudConfigRegistry
	if err := json.Unmarshal(data, &configRegistry); err != nil {
		return nil, err
	}

	return &configRegistry, nil
}

func (m *CloudManagerImpl) GetChecksum(key string) (*[32]byte, error) {
	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return nil, err
	}
	for keyInRegistry, checksum := range cloudRegistry.Configs {
		if keyInRegistry == key {
			return &checksum, nil
		}
	}
	return nil, ErrKeyNotFoundInCloud
}

func (m *CloudManagerImpl) GetAllKeys() ([]string, error) {
	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return nil, err
	}
	var keys []string
	for key := range cloudRegistry.Configs {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *CloudManagerImpl) SaveCloudConfigRegistry(cloudConfigRegistry CloudConfigRegistry) error {
	data, err := json.Marshal(cloudConfigRegistry)
	if err != nil {
		return err
	}
	return m.Provider.Upload(cloudManagerFileName, data)
}

func (m *CloudManagerImpl) SetConfig(key string, checksum [32]byte) error {
	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return err
	}
	cloudRegistry.Configs[key] = checksum
	return m.SaveCloudConfigRegistry(*cloudRegistry)
}

func (m *CloudManagerImpl) RemoveConfig(key string) error {
	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return err
	}
	for keyInRegistry := range cloudRegistry.Configs {
		if keyInRegistry == key {
			delete(cloudRegistry.Configs, key)
		}
	}

	return ErrKeyNotFoundInCloud
}

func (m *CloudManagerImpl) SetChecksum(key string, checksum [32]byte) error {
	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return err
	}
	cloudRegistry.Configs[key] = checksum
	return nil
}

func (m *CloudManagerImpl) UpdateConfig(configObj ConfigObj) error {
	checksum := sha256.Sum256(configObj.Content)

	if err := m.SetChecksum(configObj.ConfigKeyName, checksum); err != nil {
		return err
	}

	if err := m.Provider.Upload(configObj.ConfigKeyName+".json", configObj.Content); err != nil {
		return err
	}

	return nil
}

func (m *CloudManagerImpl) DownloadConfig(key string) (*ConfigObj, error) {
	data, err := m.Provider.Download(key + ".json")
	if err != nil {
		return nil, err
	}
	var configObj ConfigObj
	if err := json.Unmarshal(data, &configObj); err != nil {
		return nil, err
	}

	return &configObj, nil
}
