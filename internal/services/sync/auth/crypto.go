package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	tokenFileVersion = 1
	tokenFileAAD     = "cm-tokens-v1"
	kdfArgon2id      = "argon2id"
	saltLen          = 16
	nonceLen         = 12
	keyLen           = 32
)

var defaultKDFParams = kdfParams{
	Time:      4,
	MemoryKiB: 128 * 1024,
	Threads:   4,
}

type kdfParams struct {
	Time      uint32 `json:"time"`
	MemoryKiB uint32 `json:"memory_kib"`
	Threads   uint8  `json:"threads"`
}

type encryptedBlob struct {
	Version    int       `json:"version"`
	KDF        string    `json:"kdf"`
	KDFParams  kdfParams `json:"kdf_params"`
	Salt       string    `json:"salt"`
	Nonce      string    `json:"nonce"`
	Ciphertext string    `json:"ciphertext"`
}

func deriveKey(passphrase, salt []byte, params kdfParams) []byte {
	return argon2.IDKey(passphrase, salt, params.Time, params.MemoryKiB, params.Threads, keyLen)
}

func seal(plaintext, passphrase []byte) ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}
	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	key := deriveKey(passphrase, salt, defaultKDFParams)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ct := gcm.Seal(nil, nonce, plaintext, []byte(tokenFileAAD))

	blob := encryptedBlob{
		Version:    tokenFileVersion,
		KDF:        kdfArgon2id,
		KDFParams:  defaultKDFParams,
		Salt:       base64.StdEncoding.EncodeToString(salt),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ct),
	}
	return json.MarshalIndent(blob, "", "  ")
}

func openBlob(blob, passphrase []byte) ([]byte, error) {
	var b encryptedBlob
	if err := json.Unmarshal(blob, &b); err != nil {
		return nil, fmt.Errorf("parse token blob: %w", err)
	}
	if b.Version != tokenFileVersion {
		return nil, fmt.Errorf("unsupported token file version %d", b.Version)
	}
	if b.KDF != kdfArgon2id {
		return nil, fmt.Errorf("unsupported kdf %q", b.KDF)
	}
	salt, err := base64.StdEncoding.DecodeString(b.Salt)
	if err != nil {
		return nil, fmt.Errorf("decode salt: %w", err)
	}
	nonce, err := base64.StdEncoding.DecodeString(b.Nonce)
	if err != nil {
		return nil, fmt.Errorf("decode nonce: %w", err)
	}
	ct, err := base64.StdEncoding.DecodeString(b.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}

	key := deriveKey(passphrase, salt, b.KDFParams)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	pt, err := gcm.Open(nil, nonce, ct, []byte(tokenFileAAD))
	if err != nil {
		return nil, ErrRetrieveTokenFromStorage
	}
	return pt, nil
}
