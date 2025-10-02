package sync_services

import (
	"crypto/sha256"
	"sync"
)

type SyncService interface {
	// Authorization
	Auth(provider string) error
	Logout(provider string) error // blank provider param => logout for everyone

	// Pulling
	PullAll() ([]SyncResult, error)
	PullOne(key string) SyncResult

	// Pushing
	Push(configs []*ConfigObj, force bool) ([]*SyncResult, error)
}

type SyncServiceImpl struct {
	CloudManager CloudManager
	AuthManager  AuthManager
}

type SyncResult struct {
	ConfigObj *ConfigObj
	Error     error
}

func (s *SyncServiceImpl) Auth(provider string) error {
	return s.AuthManager.Authenticate(provider)
}

func (s *SyncServiceImpl) Logout(provider string) error {
	if provider == "" {
		return s.AuthManager.RemoveAllTokens()
	}
	return s.AuthManager.RemoveToken(provider)
}

func (s *SyncServiceImpl) PullOne(key string) SyncResult {
	configRegistry, err := s.CloudManager.GetCloudInfo()
	if err != nil {
		return SyncResult{
			ConfigObj: nil,
			Error:     err,
		}
	}
	if !configRegistry.KeyExist(key) {
		return SyncResult{
			ConfigObj: nil,
			Error:     ErrKeyNotFoundInCloud,
		}
	}

	cfgObj, err := s.CloudManager.DownloadConfig(key)
	return SyncResult{ConfigObj: cfgObj, Error: err}
}

func (s *SyncServiceImpl) PullAll() ([]SyncResult, error) {
	configRegistry, err := s.CloudManager.GetCloudInfo()
	if err != nil {
		return nil, err
	}
	keys := configRegistry.GetAllKeys()
	results := []SyncResult{}
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
	for res := range resChan {
		results = append(results, res)
	}
	return results, nil
}
func (s *SyncServiceImpl) Push(configs []*ConfigObj, force bool) ([]*SyncResult, error) {
	cloudConfigRegistry, err := s.CloudManager.GetCloudInfo()

	if err != nil {
		return nil, err
	}

	filteredConfigs := []*ConfigObj{}

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
			if cloudChecksum != localChecksum {
				cloudConfigRegistry.SetChecksum(cfg.KeyName, localChecksum)
				filteredConfigs = append(filteredConfigs, cfg)
			}
		}
	

	// Updating cloud registry after removing extra configs
	if (!cloudRegistryChanged) && len(filteredConfigs) == 0 {
		return nil, ErrNothingToPush
	}

	if err := s.CloudManager.SaveCloudConfigRegistry(*cloudConfigRegistry); err != nil {
		return nil, err
	}

	// Starting Updaing all of the colected configs
	return s.CloudManager.ConcurrentUpdateConfigs(filteredConfigs)
}

func NewSyncServiceImpl(authManager AuthManager) *SyncServiceImpl {
	token, err := authManager.GetToken("dropbox")
	var cloud CloudManager = NoopCloudManager{
		Error: err,
	}
	if err == nil {
		cloud = NewCloudManagerImpl(token)
	}
	return &SyncServiceImpl{
		AuthManager:  authManager,
		CloudManager: cloud,
	}
}
