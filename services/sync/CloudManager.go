package sync_services

import (
	"crypto/sha256"
	"encoding/json"
)

const (
	cloudManagerFileName = "cloud_manager.json"
)

type CloudManager interface {
	GetCloudInfo() (*CloudConfigRegistry, error)

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
		defaultRegistry := CloudConfigRegistry{
			Configs: make(map[string][32]byte),
		}
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
	if configRegistry.Configs == nil {
		configRegistry.Configs = make(map[string][32]byte)
	}

	return &configRegistry, nil
}

func (m *CloudManagerImpl) SaveCloudConfigRegistry(cloudConfigRegistry CloudConfigRegistry) error {
	data, err := json.Marshal(cloudConfigRegistry)
	if err != nil {
		return err
	}
	return m.Provider.Upload(cloudManagerFileName, data)
}

func (m *CloudManagerImpl) UpdateConfig(configObj ConfigObj) error {
	checksum := sha256.Sum256(configObj.Content)

	cloudRegistry, err := m.GetCloudInfo()
	if err != nil {
		return err
	}
	cloudRegistry.SetChecksum(configObj.KeyName, checksum)
	if err := m.SaveCloudConfigRegistry(*cloudRegistry); err != nil {
		return err
	}

	data, err := json.Marshal(configObj)

	if err != nil {
		return err
	}

	if err := m.Provider.Upload(configObj.KeyName+".json", data); err != nil {
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
