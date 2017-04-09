// go-bindata (https://github.com/jteeuwen/go-bindata) stub
// so beehive works also without embedded assets.
package api

import "io/ioutil"

func Asset(asset string) ([]byte, error) {
	return ioutil.ReadFile(asset)
}
