package cloud

import (
	"encoding/json"
	"sync"

	"github.com/hurtki/configsManager/internal/domain"
	sync_services "github.com/hurtki/configsManager/internal/services/sync"
)

const (
	// File that is being stored in cloud, contains all information about config's keys and their checksum
	cloudManagerFileName = "cloud_manager.json"
)

type CloudManager struct {
	Provider Provider
}

func NewCloudManager(token string) *CloudManager {
	return &CloudManager{
		Provider: NewDropboxProvider(token),
	}
}
func (m *CloudManager) GetCloudInfo() (*domain.CloudConfigRegistry, error) {
	data, err := m.Provider.Download(cloudManagerFileName)
	if err == ErrFileDoesntExist {
		defaultRegistry := domain.CloudConfigRegistry{
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

	var configRegistry domain.CloudConfigRegistry
	if err := json.Unmarshal(data, &configRegistry); err != nil {
		return nil, err
	}
	if configRegistry.Configs == nil {
		configRegistry.Configs = make(map[string][32]byte)
	}

	return &configRegistry, nil
}

func (m *CloudManager) SaveCloudConfigRegistry(cloudConfigRegistry domain.CloudConfigRegistry) error {
	data, err := json.Marshal(cloudConfigRegistry)
	if err != nil {
		return err
	}
	return m.Provider.Upload(cloudManagerFileName, data)
}

func (m *CloudManager) UpdateConfig(configObj domain.ConfigObj) error {

	data, err := json.Marshal(configObj)

	if err != nil {
		return err
	}

	if err := m.Provider.Upload(configObj.KeyName+".json", data); err != nil {
		return err
	}

	return nil
}

func (m CloudManager) ConcurrentUpdateConfigs(configs []*domain.ConfigObj) ([]*sync_services.SyncResult, error) {

	results := []*sync_services.SyncResult{}
	resChan := make(chan sync_services.SyncResult)
	wg := &sync.WaitGroup{}
	wg.Add(len(configs))

	for _, cfg := range configs {
		go func(cfg *domain.ConfigObj) {
			defer wg.Done()
			data, err := json.Marshal(cfg)

			if err != nil {
				resChan <- sync_services.SyncResult{
					ConfigObj: cfg,
					Error:     err,
				}
				return
			}

			resChan <- sync_services.SyncResult{
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

func (m *CloudManager) DownloadConfig(key string) (*domain.ConfigObj, error) {
	data, err := m.Provider.Download(key + ".json")
	if err != nil {
		return nil, err
	}
	var configObj domain.ConfigObj
	if err := json.Unmarshal(data, &configObj); err != nil {
		return nil, err
	}

	return &configObj, nil
}
