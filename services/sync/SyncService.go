package sync_services

import (
	"fmt"
	"sync"
)

type SyncService interface {
	// Authorization
	Auth(provider string, token string) error
	Logout(provider string) error // blank provider param => logout for everyone

	// Pulling
	PullAll() []SyncResult
	PullOne(key string) SyncResult

	// Pushing
	Push(configs []*ConfigObj, force bool) []SyncResult
}

type SyncServiceImpl struct {
	CloudManager CloudManager
	AuthManager  AuthManager
}

type SyncResult struct {
	ConfigObj *ConfigObj
	Error     error
}

func (s *SyncServiceImpl) Auth(provider, token string) error {
	if token == "" {
		return s.AuthManager.Authenticate(provider)
	} else {
		return fmt.Errorf("token authorisation not supported")
	}
}

func (s *SyncServiceImpl) Logout(provider string) error {
	if provider == "" {
		return s.AuthManager.RemoveAllTokens()
	}
	return s.AuthManager.RemoveToken(provider)
}

func (s *SyncServiceImpl) PullOne(key string) SyncResult {
	cfgObj, err := s.CloudManager.DownloadConfig(key)
	return SyncResult{ConfigObj: cfgObj, Error: err}
}

func (s *SyncServiceImpl) PullAll() []SyncResult {
	configRegistry, err := s.CloudManager.GetCloudInfo()
	if err != nil {
		return nil
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
	return results
}

func (s *SyncServiceImpl) Push(configs []*ConfigObj, force bool) []SyncResult {
	// WIP, need to check checksums and updating checksums not one by one

	results := []SyncResult{}
	resChan := make(chan SyncResult)
	wg := &sync.WaitGroup{}
	wg.Add(len(configs))

	for _, cfg := range configs {
		go func(cfg *ConfigObj) {
			defer wg.Done()
			err := s.CloudManager.UpdateConfig(*cfg)
			resChan <- SyncResult{
				ConfigObj: cfg,
				Error:     err,
			}
		}(cfg)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		results = append(results, res)
	}
	return results
}

func NewSyncServiceImpl(authManager AuthManager) *SyncServiceImpl {
	token, err := authManager.GetToken("dropbox")
	var cloud CloudManager
	if err != nil {
		cloud = NoopCloudManager{}
	} else {
		cloud = NewCloudManagerImpl(token)
	}

	return &SyncServiceImpl{
		AuthManager:  authManager,
		CloudManager: cloud,
	}
}
