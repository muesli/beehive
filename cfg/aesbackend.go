/*
 *    Copyright (C) 2020 Sergio Rubio
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Sergio Rubio <sergio@rubio.im>
 */

package cfg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
)

// PasswordEnvVar defines the environment variable name that should
// contain the configuration password.
const PasswordEnvVar = "BEEHIVE_CONFIG_PASSWORD"

// EncryptedHeaderPrefix is added to the encrypted configuration
// to make it possible to detect it's an encrypted configuration file
const EncryptedHeaderPrefix = "beehiveconf+"

// AESBackend symmetrically encrypts the configuration file using AES-GCM
type AESBackend struct{}

// NewAESBackend creates the backend.
//
// Given the password is required to encrypt/decrypt the configuration, if the
// URL passed doesn't have a password or PasswordEnvVar is not defined,
// it'll return an error.
func NewAESBackend(u *url.URL) (*AESBackend, error) {
	if _, err := getPassword(u); err != nil {
		return nil, err
	}

	return &AESBackend{}, nil
}

// IsEncrypted returns true and no error if the configuration is encrypted
//
// If the error returned is not nil, an error was returned while opening or
// reading the file.
func IsEncrypted(u *url.URL) (bool, error) {
	fixWinURL(u)
	f, err := os.Open(u.Path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	b := make([]byte, 12)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}

	if string(b) != EncryptedHeaderPrefix {
		return false, nil
	}

	return true, nil
}

// Load configuration file from the given URL and decrypt it
func (b *AESBackend) Load(u *url.URL) (*Config, error) {
	fixWinURL(u)
	config := &Config{url: u}

	if !exist(u.Path) {
		return config, nil
	}

	ciphertext, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return nil, err
	}
	ftype := ciphertext[0:12]
	if string(ftype) != EncryptedHeaderPrefix {
		return nil, errors.New("encrypted configuration header not valid")
	}

	p, err := getPassword(u)
	if err != nil {
		return nil, err
	}

	plaintext, err := decrypt(ciphertext[12:], []byte(p))
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
	fixWinURL(u)
	cfgDir := filepath.Dir(u.Path)
	if !exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	p, err := getPassword(config.URL())
	if err != nil {
		return err
	}

	ciphertext, err := encrypt(j, []byte(p))
	if err != nil {
		return err
	}

	marked := []byte(EncryptedHeaderPrefix)
	err = ioutil.WriteFile(u.Path, append(marked, ciphertext...), 0644)

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

func getPassword(u *url.URL) (string, error) {
	p := os.Getenv(PasswordEnvVar)
	if p != "" {
		return p, nil
	}

	p = u.User.Username()
	if p != "" {
		return p, nil
	}

	return "", errors.New("password to encrypt or decrypt the config file not available")
}
