package cfg

import (
	"net/url"
	"path/filepath"
	"runtime"
)

// URL wraps net/url.URL to deal with some platform specific issues when
// parsing configuration urls.
type URL struct {
	nurl   *url.URL
	Path   string
	Host   string
	Scheme string
	User   *url.Userinfo
	Raw    string
}

// ParseURL mimicks net/url URL.Parse
func ParseURL(rawurl string) (*URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	fixWinURL(u)
	curl := URL{}
	curl.nurl = u
	curl.User = u.User
	curl.Path = u.Path
	curl.Host = u.Host
	curl.Scheme = u.Scheme
	curl.Raw = rawurl

	return &curl, nil
}

func (u *URL) NetURL() *url.URL {
	return u.nurl
}

// Workaround Go URL parsing in Windows, where URL.Host contains the drive
// info for a URL like file://c:/path/to/beehive.config
//
// no-op in non-Windows OSes
func fixWinURL(u *url.URL) {
	if runtime.GOOS == "windows" {
		u.Path = filepath.Join(u.Host, u.Path)
		u.Host = ""
	}
}
