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
		err := s.AuthManager.Authenticate(provider)
		if err != nil {
			return err
		}
		token, err = s.AuthManager.GetToken(provider)
		if err != nil {
			return err
		}
		s.CloudManager = NewCloudManagerImpl(token)
		return nil
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
	keys, err := s.CloudManager.GetAllKeys()
	if err != nil {
		return nil
	}
	results := make([]SyncResult, len(keys))
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
	results := make([]SyncResult, len(configs))
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

func NewSyncServiceImpl(AuthManager AuthManager) *SyncServiceImpl {
	return &SyncServiceImpl{
		AuthManager:  AuthManager,
		CloudManager: NoopCloudManager{},
	}
}
