package sync_services

import (
	"github.com/99designs/keyring"
)

var (
	keyringServiceName         = "configsManager"
	dropboxAccessTokenKeyName  = "cm_dropbox_access_token"
	dropboxRefreshTokenKeyName = "cm_dropbox_access_token"
)

var (
	allKeys = []string{
		dropboxAccessTokenKeyName,
		dropboxRefreshTokenKeyName,
	}
)

type TokenStore interface {
	SaveToken(providerName string, tokenPair TokenPair) error
	LoadToken(providerName string) (*TokenPair, error)
	DeleteToken(providerName string) error
}

type TokenStoreImpl struct {
	ring keyring.Keyring
}

type TokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (s *TokenStoreImpl) SaveToken(providerName string, tokenPair TokenPair) error {

	if providerName == "dropbox" {
		if err := s.ring.Set(keyring.Item{
			Key:  dropboxAccessTokenKeyName,
			Data: []byte(tokenPair.Access),
		}); err != nil {
			return err
		}
		if err := s.ring.Set(keyring.Item{
			Key:  dropboxRefreshTokenKeyName,
			Data: []byte(tokenPair.Refresh),
		}); err != nil {
			return err
		}
	} else {
		return ErrAuthProviderDoesntExist
	}

	return nil

}

func (s *TokenStoreImpl) LoadToken(providerName string) (*TokenPair, error) {

	if providerName == "dropbox" {
		accessToken, err := s.ring.Get(dropboxAccessTokenKeyName)
		if err != nil {
			if err == keyring.ErrKeyNotFound {
				return nil, ErrTokenNotFoundInSecrets
			}
			return nil, err
		}
		refreshToken, err := s.ring.Get(dropboxRefreshTokenKeyName)
		if err != nil {
			if err == keyring.ErrKeyNotFound {
				return nil, ErrTokenNotFoundInSecrets
			}
			return nil, err
		}
		return &TokenPair{
			Access:  string(accessToken.Data),
			Refresh: string(refreshToken.Data),
		}, nil
	} else {
		return nil, ErrAuthProviderDoesntExist
	}
}

func (s *TokenStoreImpl) DeleteToken(providerName string) error {
	switch providerName {
	case "dropbox":
		return s.ring.Remove(dropboxAccessTokenKeyName)
	case "":
		for _, key := range allKeys {
			err := s.ring.Remove(key)
			if err != nil {
				return err
			}
		}
	default:
		return ErrAuthProviderDoesntExist
	}
	return nil
}

func NewTokenStoreImpl() *TokenStoreImpl {
	ring, _ := keyring.Open(keyring.Config{
		ServiceName: keyringServiceName,
	})
	return &TokenStoreImpl{
		ring: ring,
	}
}
