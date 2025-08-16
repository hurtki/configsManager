package sync_services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	dropboxAppKey = "6j9zzv3szf21pid"
)

type AuthResult struct {
	Success bool
}

type AuthManager interface {
	Authenticate(providerName string) error
	RefreshToken(providerName string) error
	GetToken(providerName string) (string, error)
	RemoveToken(providerName string) error
	RemoveAllTokens() error
}

type AuthManagerImpl struct {
	TokenStore TokenStore
}

func (m *AuthManagerImpl) randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (m *AuthManagerImpl) authenticateDropbox() error {
	codeVerifier := m.randomString(rand.Intn(80) + 43)
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	fmt.Printf("https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code&code_challenge=%s&code_challenge_method=S256&token_access_type=offline\n", dropboxAppKey, codeChallenge)
	var tokenFromUser string
	fmt.Print("Paste code you got from dropbox:")
	if _, err := fmt.Scanln(&tokenFromUser); err != nil {
		return err
	}

	reqUrl := "https://api.dropboxapi.com/oauth2/token"

	reqBody := url.Values{}
	reqBody.Set("code", tokenFromUser)
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("client_id", dropboxAppKey)
	reqBody.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest(
		"POST",
		reqUrl,
		strings.NewReader(reqBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response body")
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parsing into access and refresh tokens
	var tokenPair TokenPair
	if err := json.Unmarshal(data, &tokenPair); err != nil {
		return err
	}

	m.TokenStore.SaveToken("dropbox", tokenPair)

	return nil
}

func (m *AuthManagerImpl) refreshDropbox() error {
	tokenPair, err := m.TokenStore.LoadToken("dropbox")
	if err != nil {
		return err
	}
	reqBody := url.Values{}
    reqBody.Set("grant_type", "refresh_token")
    reqBody.Set("refresh_token", tokenPair.Refresh)
    reqBody.Set("client_id", dropboxAppKey)
    

    req, err := http.NewRequest("POST", "https://api.dropboxapi.com/oauth2/token", strings.NewReader(reqBody.Encode()))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // parsing only access_token field
    var respJson struct {
        Access string `json:"access_token"`
    }
    if err := json.Unmarshal(data, &respJson); err != nil {
        return err
    }

    tokenPair.Access = respJson.Access // updating only access field in tokenPair
    return m.TokenStore.SaveToken("dropbox", *tokenPair)
}

func (m *AuthManagerImpl) getTokenDropbox() (string, error) {

	tokenPair, err := m.TokenStore.LoadToken("dropbox")
	if err != nil {
		return "", err
	}
	return tokenPair.Access, nil
}

func (m *AuthManagerImpl) Authenticate(providerName string) error {
	if providerName == "dropbox" {
		return m.authenticateDropbox()
	}
	return ErrAuthProviderDoesntExist
}

func (m *AuthManagerImpl) RefreshToken(providerName string) error {
	if providerName == "dropbox" {
		return m.refreshDropbox()
	}
	return ErrAuthProviderDoesntExist
}

func (m *AuthManagerImpl) GetToken(providerName string) (string, error) {
	if providerName == "dropbox" {
		return m.getTokenDropbox()
	}
	return "", ErrAuthProviderDoesntExist
}

func (m *AuthManagerImpl) RemoveToken(providerName string) error {
	return m.TokenStore.DeleteToken(providerName)
}

func (m *AuthManagerImpl) RemoveAllTokens() error {
	return m.TokenStore.DeleteToken("dropbox")
}

func NewAuthManagerImpl() AuthManager {
	return &AuthManagerImpl{
		TokenStore: NewTokenStoreImpl(),
	}
}
