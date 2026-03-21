package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

var (
	dropboxAppKey = "6j9zzv3szf21pid"
)

type AuthManager struct {
	TokenStore TokenStore
}

func NewAuthManager(TokenStore TokenStore) *AuthManager {
	return &AuthManager{
		TokenStore: TokenStore,
	}
}

func (m *AuthManager) randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (m *AuthManager) authenticateDropbox() error {
	codeVerifier := m.randomString(rand.Intn(80) + 43)
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

	fmt.Printf("\x1b]8;;https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code&code_challenge=%s&code_challenge_method=S256&token_access_type=offline\x1b\\Click here to authorize\x1b]8;;\x1b\\\n", dropboxAppKey, codeChallenge)

	fmt.Print("Paste code you got from Dropbox: ")
	codeBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()
	tokenFromUser := string(codeBytes)

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

	if resp.StatusCode != 200 {
		return ErrInvalidDropboxCode
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parsing into access and refresh tokens
	var tokenPair TokenPair
	if err := json.Unmarshal(data, &tokenPair); err != nil {
		return err
	}

	return m.TokenStore.SaveToken("dropbox", tokenPair)
}

func (m *AuthManager) getTokenDropbox() (string, error) {

	tokenPair, err := m.TokenStore.LoadToken("dropbox")
	if err != nil {
		return "", err
	}
	return tokenPair.Access, nil
}

func (m *AuthManager) Authenticate(providerName string) error {
	if providerName == "dropbox" {
		return m.authenticateDropbox()
	}
	return ErrAuthProviderDoesntExist
}

func (m *AuthManager) GetToken(providerName string) (string, error) {
	if providerName == "dropbox" {
		return m.getTokenDropbox()
	}
	return "", ErrAuthProviderDoesntExist
}

func (m *AuthManager) RemoveToken(providerName string) error {
	return m.TokenStore.DeleteToken(providerName)
}

func (m *AuthManager) RemoveAllTokens() error {
	return m.TokenStore.DeleteToken("dropbox")
}

// =================== WIP ====================
// ============================================
func (m *AuthManager) RefreshToken(providerName string) error {
	if providerName == "dropbox" {
		return m.refreshDropbox()
	}
	return ErrAuthProviderDoesntExist
}

func (m *AuthManager) refreshDropbox() error {
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response")
		}
	}()
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
