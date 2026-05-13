package auth

import (
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
