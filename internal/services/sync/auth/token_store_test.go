package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fakeOs struct{ home string }

func (f *fakeOs) GetHomeDir() (string, error) { return f.home, nil }

func newStore(t *testing.T, pass string) *TokenStoreImpl {
	t.Helper()
	store := NewTokenStoreImpl(&fakeOs{home: t.TempDir()})
	store.promptFn = func(_ string, _ bool) (string, error) { return pass, nil }
	return store
}

func TestTokenStore_SaveLoadRoundtrip(t *testing.T) {
	store := newStore(t, "correct horse battery staple")

	want := TokenPair{Access: "acc-123", Refresh: "ref-456"}
	if err := store.SaveToken("dropbox", want); err != nil {
		t.Fatalf("SaveToken: %v", err)
	}

	got, err := store.LoadToken("dropbox")
	if err != nil {
		t.Fatalf("LoadToken: %v", err)
	}
	if got.Access != want.Access || got.Refresh != want.Refresh {
		t.Fatalf("roundtrip mismatch: got %+v, want %+v", got, want)
	}
}

func TestTokenStore_PersistsAcrossInstances(t *testing.T) {
	home := t.TempDir()
	pass := "p4ssw0rd"

	s1 := NewTokenStoreImpl(&fakeOs{home: home})
	s1.promptFn = func(_ string, _ bool) (string, error) { return pass, nil }
	if err := s1.SaveToken("dropbox", TokenPair{Access: "a", Refresh: "r"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	s2 := NewTokenStoreImpl(&fakeOs{home: home})
	s2.promptFn = func(_ string, _ bool) (string, error) { return pass, nil }
	got, err := s2.LoadToken("dropbox")
	if err != nil {
		t.Fatalf("load on fresh instance: %v", err)
	}
	if got.Access != "a" || got.Refresh != "r" {
		t.Fatalf("got %+v", got)
	}
}

func TestTokenStore_WrongPassphrase(t *testing.T) {
	home := t.TempDir()

	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "right", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "x", Refresh: "y"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	loader := NewTokenStoreImpl(&fakeOs{home: home})
	loader.promptFn = func(_ string, _ bool) (string, error) { return "wrong", nil }
	_, err := loader.LoadToken("dropbox")
	if !errors.Is(err, ErrRetrieveTokenFromStorage) {
		t.Fatalf("expected ErrRetrieveTokenFromStorage, got %v", err)
	}
}

func TestTokenStore_FileMissing(t *testing.T) {
	store := newStore(t, "anything")

	_, err := store.LoadToken("dropbox")
	if !errors.Is(err, ErrTokenNotFoundInSecrets) {
		t.Fatalf("expected ErrTokenNotFoundInSecrets, got %v", err)
	}
}

func TestTokenStore_FilePermissions(t *testing.T) {
	home := t.TempDir()
	store := NewTokenStoreImpl(&fakeOs{home: home})
	store.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }

	if err := store.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if mode := info.Mode().Perm(); mode != 0600 {
		t.Fatalf("expected 0600 perms, got %o", mode)
	}
}

func TestTokenStore_DeleteProviderRemovesFile(t *testing.T) {
	home := t.TempDir()
	store := NewTokenStoreImpl(&fakeOs{home: home})
	store.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }

	if err := store.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := store.DeleteToken("dropbox"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file to be removed, stat err: %v", err)
	}
}

func TestTokenStore_DeleteAll(t *testing.T) {
	home := t.TempDir()
	store := NewTokenStoreImpl(&fakeOs{home: home})
	store.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }

	if err := store.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := store.DeleteToken(""); err != nil {
		t.Fatalf("delete-all: %v", err)
	}
	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file to be removed: %v", err)
	}
}

func TestTokenStore_UnknownProvider(t *testing.T) {
	store := newStore(t, "pw")
	if err := store.SaveToken("googledrive", TokenPair{}); !errors.Is(err, ErrAuthProviderDoesntExist) {
		t.Fatalf("save: expected ErrAuthProviderDoesntExist, got %v", err)
	}
	if _, err := store.LoadToken("googledrive"); !errors.Is(err, ErrAuthProviderDoesntExist) {
		t.Fatalf("load: expected ErrAuthProviderDoesntExist, got %v", err)
	}
}

func readBlobSalt(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read blob: %v", err)
	}
	var b struct {
		Salt string `json:"salt"`
	}
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("parse blob: %v", err)
	}
	return b.Salt
}

func TestTokenStore_SaltStableAcrossSaves(t *testing.T) {
	home := t.TempDir()
	store := NewTokenStoreImpl(&fakeOs{home: home})
	store.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }

	if err := store.SaveToken("dropbox", TokenPair{Access: "a", Refresh: "r"}); err != nil {
		t.Fatalf("first save: %v", err)
	}
	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	salt1 := readBlobSalt(t, path)

	if err := store.SaveToken("dropbox", TokenPair{Access: "a2", Refresh: "r2"}); err != nil {
		t.Fatalf("second save: %v", err)
	}
	salt2 := readBlobSalt(t, path)

	if salt1 != salt2 {
		t.Fatalf("salt rotated unexpectedly: %s vs %s", salt1, salt2)
	}
}

func TestTokenStore_SaltRotatesAfterDelete(t *testing.T) {
	home := t.TempDir()
	store := NewTokenStoreImpl(&fakeOs{home: home})
	store.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }

	if err := store.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	saltBefore := readBlobSalt(t, path)

	if err := store.DeleteToken(""); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := store.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save after delete: %v", err)
	}
	saltAfter := readBlobSalt(t, path)

	if saltBefore == saltAfter {
		t.Fatalf("expected fresh salt after delete+create, got the same: %s", saltBefore)
	}
}

func TestTokenStore_DeriveSessionKey(t *testing.T) {
	home := t.TempDir()
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a", Refresh: "r"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	deriver := NewTokenStoreImpl(&fakeOs{home: home})
	deriver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	key, err := deriver.DeriveSessionKey()
	if err != nil {
		t.Fatalf("DeriveSessionKey: %v", err)
	}
	if len(key) != keyLen {
		t.Fatalf("expected %d-byte key, got %d", keyLen, len(key))
	}
}

func TestTokenStore_DeriveSessionKeyMissingFile(t *testing.T) {
	store := newStore(t, "pw")
	_, err := store.DeriveSessionKey()
	if !errors.Is(err, ErrTokenNotFoundInSecrets) {
		t.Fatalf("expected ErrTokenNotFoundInSecrets, got %v", err)
	}
}

func TestTokenStore_DeriveSessionKeyWrongPassphrase(t *testing.T) {
	home := t.TempDir()
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "right", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	deriver := NewTokenStoreImpl(&fakeOs{home: home})
	deriver.promptFn = func(_ string, _ bool) (string, error) { return "wrong", nil }
	_, err := deriver.DeriveSessionKey()
	if !errors.Is(err, ErrRetrieveTokenFromStorage) {
		t.Fatalf("expected ErrRetrieveTokenFromStorage, got %v", err)
	}
}

func TestTokenStore_LoadWithEnvSessionKey(t *testing.T) {
	home := t.TempDir()

	// Save with passphrase + derive the session key.
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a", Refresh: "r"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	key, err := saver.DeriveSessionKey()
	if err != nil {
		t.Fatalf("derive: %v", err)
	}

	// Fresh instance: passphrase prompt must NOT be called.
	t.Setenv(sessionKeyEnv, base64.StdEncoding.EncodeToString(key))
	loader := NewTokenStoreImpl(&fakeOs{home: home})
	loader.promptFn = func(_ string, _ bool) (string, error) {
		t.Fatal("prompt should not be called when CM_SYNC_KEY is valid")
		return "", nil
	}
	got, err := loader.LoadToken("dropbox")
	if err != nil {
		t.Fatalf("LoadToken: %v", err)
	}
	if got.Access != "a" || got.Refresh != "r" {
		t.Fatalf("got %+v", got)
	}
}

func TestTokenStore_StaleEnvKeyFallsBackToPrompt(t *testing.T) {
	home := t.TempDir()
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	// Random 32-byte key — won't match the file.
	bogus := make([]byte, keyLen)
	for i := range bogus {
		bogus[i] = byte(i)
	}
	t.Setenv(sessionKeyEnv, base64.StdEncoding.EncodeToString(bogus))

	promptCalls := 0
	loader := NewTokenStoreImpl(&fakeOs{home: home})
	loader.promptFn = func(_ string, _ bool) (string, error) {
		promptCalls++
		return "pw", nil
	}
	pair, err := loader.LoadToken("dropbox")
	if err != nil {
		t.Fatalf("LoadToken: %v", err)
	}
	if pair.Access != "a" {
		t.Fatalf("got %+v", pair)
	}
	if promptCalls != 1 {
		t.Fatalf("expected exactly one prompt call, got %d", promptCalls)
	}
}

func TestTokenStore_MalformedEnvKeyFallsBackToPrompt(t *testing.T) {
	home := t.TempDir()
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	t.Setenv(sessionKeyEnv, "not-base64@@@")

	loader := NewTokenStoreImpl(&fakeOs{home: home})
	loader.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if _, err := loader.LoadToken("dropbox"); err != nil {
		t.Fatalf("LoadToken: %v", err)
	}
}

func TestTokenStore_SaveAfterEnvUnlockKeepsSalt(t *testing.T) {
	home := t.TempDir()
	saver := NewTokenStoreImpl(&fakeOs{home: home})
	saver.promptFn = func(_ string, _ bool) (string, error) { return "pw", nil }
	if err := saver.SaveToken("dropbox", TokenPair{Access: "a"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	path := filepath.Join(home, ".config", "configsManager", tokensFileName)
	saltBefore := readBlobSalt(t, path)

	key, err := saver.DeriveSessionKey()
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	t.Setenv(sessionKeyEnv, base64.StdEncoding.EncodeToString(key))

	loader := NewTokenStoreImpl(&fakeOs{home: home})
	loader.promptFn = func(_ string, _ bool) (string, error) {
		t.Fatal("prompt should not run during env-unlock flow")
		return "", nil
	}
	if err := loader.SaveToken("dropbox", TokenPair{Access: "a2"}); err != nil {
		t.Fatalf("save with env key: %v", err)
	}
	saltAfter := readBlobSalt(t, path)
	if saltBefore != saltAfter {
		t.Fatalf("salt rotated under env-unlock: %s -> %s", saltBefore, saltAfter)
	}
}
