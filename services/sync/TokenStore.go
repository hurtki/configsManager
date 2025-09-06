package sync_services

import (
	"strings"

	"github.com/99designs/keyring"
)

var (
	keyringServiceName         = "configsManager"
	dropboxAccessTokenKeyName  = "cm_dropbox_access_token"
	dropboxRefreshTokenKeyName = "cm_dropbox_refresh_token"
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
	Access  string
	Refresh string
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
			if strings.Contains(err.Error(), "integrity check failed") {
				return nil, ErrRetrieveTokenFromStorage
			}

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
		for _, key := range []string{dropboxAccessTokenKeyName, dropboxRefreshTokenKeyName} {
			if err := s.ring.Remove(key); err != nil && err != keyring.ErrKeyNotFound {
				return err
			}
		}
		return nil
	case "":
		for _, key := range allKeys {
			if err := s.ring.Remove(key); err != nil && err != keyring.ErrKeyNotFound {
				return err
			}
		}
		return nil
	default:
		return ErrAuthProviderDoesntExist
	}

}

func NewTokenStoreImpl() *TokenStoreImpl {
	ring, _ := keyring.Open(keyring.Config{
		AllowedBackends: []keyring.BackendType{
			keyring.KWalletBackend,
			keyring.SecretServiceBackend,
			keyring.KeychainBackend,
			keyring.PassBackend,
			keyring.FileBackend,
		},
		ServiceName:      keyringServiceName,
		FileDir:          "~/.config/configsManager/sync_tokens/",
		FilePasswordFunc: keyring.TerminalPrompt,
	})

	return &TokenStoreImpl{
		ring: ring,
	}
}
