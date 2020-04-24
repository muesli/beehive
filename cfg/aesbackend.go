package cfg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
)

// AESBackend symmetrically encrypts the configuration file using AES-GCM
type AESBackend struct{}

// NewAESBackend create the backend
func NewAESBackend() *AESBackend {
	return &AESBackend{}
}

// Load configuration file from the given URL and decrypt it
func (b *AESBackend) Load(u *url.URL) (*Config, error) {
	config := &Config{url: u}

	if !exist(u.Path) {
		return config, nil
	}

	ciphertext, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return nil, err
	}

	p, ok := u.User.Password()
	if !ok {
		return nil, fmt.Errorf("No password specified!")
	}

	plaintext, err := decrypt(ciphertext, []byte(p))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(plaintext, config)
	if err != nil {
		return nil, err
	}

	config.backend = b
	config.url = u

	return config, nil

}

// Save encrypts then saves the configuration
func (b *AESBackend) Save(config *Config) error {
	u := config.URL()
	cfgDir := filepath.Dir(u.Path)
	if !exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	p, ok := u.User.Password()
	if !ok {
		return fmt.Errorf("No password specified!")
	}

	ciphertext, err := encrypt(j, []byte(p))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(u.Path, ciphertext, 0644)

	return err
}

func encrypt(data, key []byte) ([]byte, error) {
	key, salt, err := deriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func decrypt(data, key []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]
	key, _, err := deriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
