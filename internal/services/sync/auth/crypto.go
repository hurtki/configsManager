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

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}

// sealWithKey encrypts plaintext using the already-derived 32-byte key and the
// provided salt (which is recorded in the blob so future decrypts can re-derive
// the same key). A fresh random nonce is generated for every call.
func sealWithKey(plaintext, key, salt []byte) ([]byte, error) {
	if len(key) != keyLen {
		return nil, fmt.Errorf("invalid key length %d", len(key))
	}
	if len(salt) != saltLen {
		return nil, fmt.Errorf("invalid salt length %d", len(salt))
	}
	nonce, err := randomBytes(nonceLen)
	if err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
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

// seal generates a fresh salt, derives a key from the passphrase and produces
// an encrypted blob. Used only when there is no existing file to reuse a salt
// from (first save or after the file was deleted).
func seal(plaintext, passphrase []byte) (blob, key, salt []byte, err error) {
	salt, err = randomBytes(saltLen)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("generate salt: %w", err)
	}
	key = deriveKey(passphrase, salt, defaultKDFParams)
	blob, err = sealWithKey(plaintext, key, salt)
	if err != nil {
		return nil, nil, nil, err
	}
	return blob, key, salt, nil
}

func parseBlob(blob []byte) (*encryptedBlob, []byte, []byte, []byte, error) {
	var b encryptedBlob
	if err := json.Unmarshal(blob, &b); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parse token blob: %w", err)
	}
	if b.Version != tokenFileVersion {
		return nil, nil, nil, nil, fmt.Errorf("unsupported token file version %d", b.Version)
	}
	if b.KDF != kdfArgon2id {
		return nil, nil, nil, nil, fmt.Errorf("unsupported kdf %q", b.KDF)
	}
	salt, err := base64.StdEncoding.DecodeString(b.Salt)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("decode salt: %w", err)
	}
	nonce, err := base64.StdEncoding.DecodeString(b.Nonce)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("decode nonce: %w", err)
	}
	ct, err := base64.StdEncoding.DecodeString(b.Ciphertext)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("decode ciphertext: %w", err)
	}
	return &b, salt, nonce, ct, nil
}

func decryptWithKey(key, nonce, ct []byte) ([]byte, error) {
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

// openBlob decrypts a blob using a passphrase. Returns the plaintext, the salt
// that was stored in the blob and the derived key (so callers can cache it).
func openBlob(blob, passphrase []byte) (plaintext, salt, key []byte, err error) {
	b, salt, nonce, ct, err := parseBlob(blob)
	if err != nil {
		return nil, nil, nil, err
	}
	key = deriveKey(passphrase, salt, b.KDFParams)
	pt, err := decryptWithKey(key, nonce, ct)
	if err != nil {
		return nil, nil, nil, err
	}
	return pt, salt, key, nil
}

// openBlobWithKey decrypts a blob with an already-derived 32-byte key, skipping
// the KDF step. Used for the fast unlock path (CM_SYNC_KEY env var).
func openBlobWithKey(blob, key []byte) (plaintext, salt []byte, err error) {
	if len(key) != keyLen {
		return nil, nil, fmt.Errorf("invalid key length %d", len(key))
	}
	_, salt, nonce, ct, err := parseBlob(blob)
	if err != nil {
		return nil, nil, err
	}
	pt, err := decryptWithKey(key, nonce, ct)
	if err != nil {
		return nil, nil, err
	}
	return pt, salt, nil
}
