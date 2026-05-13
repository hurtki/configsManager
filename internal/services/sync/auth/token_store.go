package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const tokensFileName = "tokens.dat"

type TokenStore interface {
	SaveToken(providerName string, tokenPair TokenPair) error
	LoadToken(providerName string) (*TokenPair, error)
	DeleteToken(providerName string) error
}

type TokenPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type OsService interface {
	GetHomeDir() (string, error)
}

type promptFunc func(prompt string, confirm bool) (string, error)

type TokenStoreImpl struct {
	osService  OsService
	promptFn   promptFunc
	cachedPass []byte
}

type tokensPlaintext struct {
	Providers map[string]TokenPair `json:"providers"`
}

func NewTokenStoreImpl(osService OsService) *TokenStoreImpl {
	return &TokenStoreImpl{
		osService: osService,
		promptFn:  defaultPrompt,
	}
}

func (s *TokenStoreImpl) tokensPath() (string, error) {
	home, err := s.osService.GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "configsManager", tokensFileName), nil
}

func (s *TokenStoreImpl) getPassphrase(prompt string, confirm bool) ([]byte, error) {
	if s.cachedPass != nil {
		return s.cachedPass, nil
	}
	pass, err := s.promptFn(prompt, confirm)
	if err != nil {
		return nil, err
	}
	s.cachedPass = []byte(pass)
	return s.cachedPass, nil
}

func (s *TokenStoreImpl) loadPlaintext() (*tokensPlaintext, error) {
	path, err := s.tokensPath()
	if err != nil {
		return nil, err
	}
	blob, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &tokensPlaintext{Providers: map[string]TokenPair{}}, nil
	}
	if err != nil {
		return nil, err
	}

	pass, err := s.getPassphrase("Enter passphrase: ", false)
	if err != nil {
		return nil, err
	}
	pt, err := openBlob(blob, pass)
	if err != nil {
		if errors.Is(err, ErrRetrieveTokenFromStorage) {
			s.cachedPass = nil
		}
		return nil, err
	}

	var doc tokensPlaintext
	if err := json.Unmarshal(pt, &doc); err != nil {
		return nil, err
	}
	if doc.Providers == nil {
		doc.Providers = map[string]TokenPair{}
	}
	return &doc, nil
}

func (s *TokenStoreImpl) savePlaintext(doc *tokensPlaintext) error {
	path, err := s.tokensPath()
	if err != nil {
		return err
	}

	confirm := false
	if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) && s.cachedPass == nil {
		confirm = true
	}
	pass, err := s.getPassphrase("Create passphrase: ", confirm)
	if err != nil {
		return err
	}

	pt, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	blob, err := seal(pt, pass)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, blob, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *TokenStoreImpl) SaveToken(providerName string, tokenPair TokenPair) error {
	if providerName != "dropbox" {
		return ErrAuthProviderDoesntExist
	}
	doc, err := s.loadPlaintext()
	if err != nil {
		return err
	}
	doc.Providers[providerName] = tokenPair
	return s.savePlaintext(doc)
}

func (s *TokenStoreImpl) LoadToken(providerName string) (*TokenPair, error) {
	if providerName != "dropbox" {
		return nil, ErrAuthProviderDoesntExist
	}
	path, err := s.tokensPath()
	if err != nil {
		return nil, err
	}
	if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
		return nil, ErrTokenNotFoundInSecrets
	}
	doc, err := s.loadPlaintext()
	if err != nil {
		return nil, err
	}
	pair, ok := doc.Providers[providerName]
	if !ok {
		return nil, ErrTokenNotFoundInSecrets
	}
	return &pair, nil
}

func (s *TokenStoreImpl) DeleteToken(providerName string) error {
	path, err := s.tokensPath()
	if err != nil {
		return err
	}
	switch providerName {
	case "dropbox":
		if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
			return nil
		}
		doc, err := s.loadPlaintext()
		if err != nil {
			return err
		}
		if _, ok := doc.Providers[providerName]; !ok {
			return nil
		}
		delete(doc.Providers, providerName)
		if len(doc.Providers) == 0 {
			return os.Remove(path)
		}
		return s.savePlaintext(doc)
	case "":
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	default:
		return ErrAuthProviderDoesntExist
	}
}
