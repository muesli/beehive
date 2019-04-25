// +build !embed

// go-bindata (https://github.com/kevinburke/go-bindata) stub
// so beehive also works without embedded assets.
package api

import "io/ioutil"

func Asset(asset string) ([]byte, error) {
	return ioutil.ReadFile(asset)
}
