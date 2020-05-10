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
	"os"
	"path/filepath"
	"testing"
)

const testPassword = "foo"

func TestAESBackendLoad(t *testing.T) {
	u, _ := ParseURL("crypt://foo@foobar")
	backend, err := NewAESBackend(u)
	if err != nil {
		t.Error("The backend should return an error if no password specified")
	}

	_, err = backend.Load(u)
	if err != nil {
		t.Error("Loading an non-existing config file should not return an error")
	}

	// try to load the config from an absolute path using a URI
	u, err = ParseURL("crypto://" + testPassword + "@" + encryptedConfPath())
	if err != nil {
		t.Fatalf("Can't parse crypto URL: %v", err)
	}
	backend, err = NewAESBackend(u)
	if err != nil {
		t.Errorf("Error loading AES backend. %v", err)
	}
	conf, err := backend.Load(u)
	if err != nil {
		t.Errorf("Error loading config file fixture from absolute path %s. %v", u, err)
	}
	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}

	// try to load the config using the password from the environment
	os.Setenv(PasswordEnvVar, testPassword)
	u, err = ParseURL("crypto://" + encryptedConfPath())
	if err != nil {
		t.Fatalf("Can't parse crypto URL: %v", err)
	}
	backend, err = NewAESBackend(u)
	conf, err = backend.Load(u)
	if err != nil {
		t.Errorf("loading the config file using the environment password should work. %v", err)
	}

	// try to load the config with an invalid password
	os.Setenv(PasswordEnvVar, "")
	u, err = ParseURL("crypto://bar@" + encryptedConfPath())
	if err != nil {
		t.Fatalf("Can't parse crypto URL: %v", err)
	}
	backend, err = NewAESBackend(u)
	conf, err = backend.Load(u)
	if err == nil || err.Error() != "cipher: message authentication failed" {
		t.Errorf("loading the config file with an invalid password should fail. %v", err)
	}

	// environment password takes prececence
	os.Setenv(PasswordEnvVar, testPassword)
	u, err = ParseURL("crypto://bar@" + encryptedConfPath())
	if err != nil {
		t.Fatalf("Can't parse crypto URL: %v", err)
	}

	backend, err = NewAESBackend(u)
	if err != nil {
		t.Fatalf("Can't create AES backend: %v", err)
	}
	conf, err = backend.Load(u)
	if err != nil {
		t.Errorf("the password defined in %s should take precedence. %v", PasswordEnvVar, err)
	}

	u, err = ParseURL("crypto://" + testPassword + "@" + confPath())
	if err != nil {
		t.Fatalf("Can't parse crypto URL: %v", err)
	}
	conf, err = backend.Load(u)
	if err == nil || err.Error() != "encrypted configuration header not valid" {
		t.Errorf("Loading a non-encrypted config should error")
	}
}

func TestAESBackendSave(t *testing.T) {
	cwd, _ := os.Getwd()
	u, err := ParseURL(filepath.Join("crypto://"+testPassword+"@", cwd, "testdata", "beehive-crypto.conf"))
	backend, _ := NewAESBackend(u)
	c, err := backend.Load(u)
	if err != nil {
		t.Errorf("Failed to load config fixture from relative path %s: %v", u, err)
	}

	// Save the config file to a new absolute path using a URL
	p := encryptedTempConf()
	u, err = ParseURL("crypto://" + testPassword + "@" + p)
	c.SetURL(u.String())
	backend, _ = NewAESBackend(u)
	err = backend.Save(c)
	if err != nil {
		t.Errorf("failed to save the config to %s", u)
	}
	if !exist(p) {
		t.Errorf("configuration file wasn't saved to %s", p)
	}

	ok, err := IsEncrypted(u)
	if !ok {
		t.Errorf("encrypted config header not added. %v", err)
	}
}
