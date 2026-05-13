package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	tokensFileName = "tokens.dat"
	sessionKeyEnv  = "CM_SYNC_KEY"
)

type TokenStore interface {
	SaveToken(providerName string, tokenPair TokenPair) error
	LoadToken(providerName string) (*TokenPair, error)
	DeleteToken(providerName string) error
	DeriveSessionKey() ([]byte, error)
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
	osService   OsService
	promptFn    promptFunc
	cachedKey   []byte
	cachedSalt  []byte
	envKeyTried bool
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

func (s *TokenStoreImpl) tryEnvKey(blob []byte) ([]byte, []byte, bool) {
	if s.envKeyTried {
		return nil, nil, false
	}
	s.envKeyTried = true

	raw := os.Getenv(sessionKeyEnv)
	if raw == "" {
		return nil, nil, false
	}
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil || len(key) != keyLen {
		fmt.Fprintf(os.Stderr, "warning: %s is set but malformed; ignoring\n", sessionKeyEnv)
		return nil, nil, false
	}
	pt, salt, err := openBlobWithKey(blob, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %s did not match this token file; falling back to passphrase. Re-run `eval \"$(cm sync unlock)\"` after entering it.\n", sessionKeyEnv)
		return nil, nil, false
	}
	s.cachedKey = key
	s.cachedSalt = salt
	return pt, key, true
}

func (s *TokenStoreImpl) decryptWithCachedKey(blob []byte) ([]byte, bool) {
	if s.cachedKey == nil {
		return nil, false
	}
	pt, _, err := openBlobWithKey(blob, s.cachedKey)
	if err != nil {
		s.cachedKey = nil
		s.cachedSalt = nil
		return nil, false
	}
	return pt, true
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

	if pt, _, ok := s.tryEnvKey(blob); ok {
		return parsePlaintext(pt)
	}
	if pt, ok := s.decryptWithCachedKey(blob); ok {
		return parsePlaintext(pt)
	}

	pass, err := s.promptFn("Enter passphrase: ", false)
	if err != nil {
		return nil, err
	}
	pt, salt, key, err := openBlob(blob, []byte(pass))
	if err != nil {
		return nil, err
	}
	s.cachedKey = key
	s.cachedSalt = salt
	return parsePlaintext(pt)
}

func parsePlaintext(pt []byte) (*tokensPlaintext, error) {
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

	pt, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	var blob []byte
	switch {
	case s.cachedKey != nil && s.cachedSalt != nil:
		blob, err = sealWithKey(pt, s.cachedKey, s.cachedSalt)
		if err != nil {
			return err
		}
	default:
		// No cache → file should not exist (Load was called before any Save in
		// real flows). Treat as first-time creation: prompt with confirm.
		pass, err := s.promptFn("Create passphrase: ", true)
		if err != nil {
			return err
		}
		newBlob, key, salt, err := seal(pt, []byte(pass))
		if err != nil {
			return err
		}
		blob = newBlob
		s.cachedKey = key
		s.cachedSalt = salt
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
			s.cachedKey = nil
			s.cachedSalt = nil
			return os.Remove(path)
		}
		return s.savePlaintext(doc)
	case "":
		s.cachedKey = nil
		s.cachedSalt = nil
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	default:
		return ErrAuthProviderDoesntExist
	}
}

func (s *TokenStoreImpl) DeriveSessionKey() ([]byte, error) {
	path, err := s.tokensPath()
	if err != nil {
		return nil, err
	}
	blob, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrTokenNotFoundInSecrets
	}
	if err != nil {
		return nil, err
	}

	pass, err := s.promptFn("Enter passphrase: ", false)
	if err != nil {
		return nil, err
	}
	_, salt, key, err := openBlob(blob, []byte(pass))
	if err != nil {
		return nil, err
	}
	s.cachedKey = key
	s.cachedSalt = salt
	return key, nil
}
