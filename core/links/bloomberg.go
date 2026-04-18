package links

import (
	"io"
	"net/url"
	"strings"
)

// bloombergParamsToRemove are Bloomberg parameters that should be removed
var bloombergParamsToRemove = []string{
	"accessToken",
	"leadSource",
}

// isBloombergURL checks if a URL is from Bloomberg
func isBloombergURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "bloomberg.com")
}

// RemoveParamsFromBloombergURLs removes tracking parameters from Bloomberg URLs
func RemoveParamsFromBloombergURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isBloombergURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for _, param := range bloombergParamsToRemove {
			if q.Has(param) {
				q.Del(param)
				changed = true
			}
		}
		if changed {
			u.RawQuery = q.Encode()
		}
		u.Path = stripTrailingSlash(u.Path)
		return u
	})
}
