package sync_services

import (
	"encoding/json"
	"sync"
)

const (
	cloudManagerFileName = "cloud_manager.json"
)

type CloudManager interface {
	GetCloudInfo() (*CloudConfigRegistry, error)
	SaveCloudConfigRegistry(cloudConfigRegistry CloudConfigRegistry) error

	ConcurrentUpdateConfigs(configs []*ConfigObj) ([]*SyncResult, error)
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

	data, err := json.Marshal(configObj)

	if err != nil {
		return err
	}

	if err := m.Provider.Upload(configObj.KeyName+".json", data); err != nil {
		return err
	}

	return nil
}

func (m CloudManagerImpl) ConcurrentUpdateConfigs(configs []*ConfigObj) ([]*SyncResult, error) {

	results := []*SyncResult{}
	resChan := make(chan SyncResult)
	wg := &sync.WaitGroup{}
	wg.Add(len(configs))

	for _, cfg := range configs {
		go func(cfg *ConfigObj) {
			defer wg.Done()
			data, err := json.Marshal(cfg)

			if err != nil {
				resChan <- SyncResult{
					ConfigObj: cfg,
					Error:     err,
				}
				return
			}

			resChan <- SyncResult{
				ConfigObj: cfg,
				Error:     m.Provider.Upload(cfg.KeyName+".json", data),
			}

		}(cfg)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		results = append(results, &res)
	}

	return results, nil
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
