package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func encryptedConfPath() string {
	cwd, _ := os.Getwd()
	return fixWindowsPath(filepath.Join(cwd, "testdata", "beehive-crypto.conf"))
}

func confPath() string {
	cwd, _ := os.Getwd()
	return fixWindowsPath(filepath.Join(cwd, "testdata", "beehive.conf"))
}

func tmpConfPath() string {
	tmpdir, err := ioutil.TempDir("", "beehivetest")
	if err != nil {
		panic("Could not create temp directory")
	}
	return fixWindowsPath(filepath.Join(tmpdir, "testdata", "beehive.conf"))
}

func encryptedTempConf() string {
	tmpdir, err := ioutil.TempDir("", "beehivetest")
	if err != nil {
		panic("Could not create temp directory")
	}
	return fixWindowsPath(filepath.Join(tmpdir, "testdata", "beehive-crypto.conf"))
}
