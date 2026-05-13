package sync_services

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/hurtki/configsManager/internal/domain"
)

type AuthManager interface {
	Authenticate(providerName string) error
	GetToken(providerName string) (string, error)
	RemoveToken(providerName string) error
	RemoveAllTokens() error
	Unlock() (string, error)
	// WIP
	// RefreshToken(providerName string) error
}

// interface that represents entity to manage cloudConfigRegistry and update/download ConcurrentUpdateConfigs
// So the higher entity ( its dependency is this interface ) should first update configs with ConcurrentUpdateConfigs
// Then update CloudConfigRegistry with Sync results
// If you are only downloading configs, no need to change cloudConfigRegistry!!!
type CloudManager interface {
	GetCloudInfo() (*domain.CloudConfigRegistry, error)
	SaveCloudConfigRegistry(cloudConfigRegistry domain.CloudConfigRegistry) error

	ConcurrentUpdateConfigs(configs []*domain.ConfigObj) ([]*SyncResult, error)
	DownloadConfig(key string) (*domain.ConfigObj, error)
}

type SyncService struct {
	CloudManager CloudManager
	AuthManager  AuthManager
}

type SyncResult struct {
	ConfigObj *domain.ConfigObj
	Error     error
}

type PullAllResult struct {
	Configs              <-chan SyncResult
	ExpectedConfigsCount int
}

func NewSyncService(authManager AuthManager, cloudManager CloudManager) *SyncService {
	return &SyncService{
		CloudManager: cloudManager,
		AuthManager:  authManager,
	}
}

func (s *SyncService) Auth(provider string) error {
	return s.AuthManager.Authenticate(provider)
}

func (s *SyncService) Logout(provider string) error {
	if provider == "" {
		return s.AuthManager.RemoveAllTokens()
	}
	return s.AuthManager.RemoveToken(provider)
}

func (s *SyncService) Unlock() (string, error) {
	return s.AuthManager.Unlock()
}

func (s *SyncService) PullOne(key string) SyncResult {
	configRegistry, err := s.CloudManager.GetCloudInfo()
	if err != nil {
		return SyncResult{
			Error: err,
		}
	}
	if !configRegistry.KeyExist(key) {
		return SyncResult{
			Error: ErrKeyNotFoundInCloud,
		}
	}

	cfgObj, err := s.CloudManager.DownloadConfig(key)
	if err != nil {
		return SyncResult{ConfigObj: cfgObj, Error: fmt.Errorf("can't download config from cloud: %w", err)}
	} else {
		return SyncResult{ConfigObj: cfgObj}
	}
}

func (s *SyncService) PullAll() (PullAllResult, error) {
	configRegistry, err := s.CloudManager.GetCloudInfo()
	if err != nil {
		return PullAllResult{}, err
	}
	keys := configRegistry.GetAllKeys()
	resChan := make(chan SyncResult)
	wg := &sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			resChan <- s.PullOne(key)
		}(key)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()

	return PullAllResult{
		Configs:              resChan,
		ExpectedConfigsCount: len(keys),
	}, nil
}
func (s *SyncService) Push(configs []*domain.ConfigObj, force bool) ([]*SyncResult, error) {
	cloudConfigRegistry, err := s.CloudManager.GetCloudInfo()

	if err != nil {
		return nil, err
	}

	filteredConfigs := []*domain.ConfigObj{}

	localKeys := make(map[string]struct{})
	for _, cfg := range configs {
		localKeys[cfg.KeyName] = struct{}{}
	}

	// Updating cloud registry to make sure that there is no extra configs there
	var cloudRegistryChanged bool
	for key := range cloudConfigRegistry.Configs {
		if _, exists := localKeys[key]; !exists {
			cloudConfigRegistry.RemoveKey(key)
			cloudRegistryChanged = true
		}
	}

	// Cheking what configs from local are new or changed, and adding them to filtered
	for _, cfg := range configs {
		if cloudChecksum, ok := cloudConfigRegistry.Configs[cfg.KeyName]; ok {
			localChecksum := sha256.Sum256(cfg.Content)
			if cloudChecksum == localChecksum {
				continue // if same checksum, no need to push
			}
		}
		filteredConfigs = append(filteredConfigs, cfg)
	}

	// Updating cloud registry after removing extra configs
	if (!cloudRegistryChanged) && len(filteredConfigs) == 0 {
		return nil, ErrNothingToPush
	}

	// Starting Updaing all of the colected configs
	results, err := s.CloudManager.ConcurrentUpdateConfigs(filteredConfigs)

	// if we got error from pushing so we are exiting without any updates for cloudConfigRegistry
	if err != nil {
		return nil, err
	}

	// Using results from pushing configs
	// if there is no error for specific config => change its checksum ( beacause it was successfully pushed)
	// if there is an error no need to update checksum in cloud
	for i := range results {
		if results[i].Error == nil {
			cloudConfigRegistry.SetChecksum(results[i].ConfigObj.KeyName, sha256.Sum256(results[i].ConfigObj.Content))
		}
	}

	// saving a cloudConfigRegistry
	if err := s.CloudManager.SaveCloudConfigRegistry(*cloudConfigRegistry); err != nil {
		return nil, err
	}
	return results, nil
}
