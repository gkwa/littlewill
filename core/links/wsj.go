package links

import (
	"io"
	"net/url"
	"strings"
)

// wsjParamsToRemove are Wall Street Journal parameters that should be removed
var wsjParamsToRemove = []string{
	"mod",
	"reflink",
	"ref",
	"st",
}

// isWSJURL checks if a URL is from Wall Street Journal
func isWSJURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "wsj.com")
}

// RemoveParamsFromWSJURLs removes tracking parameters from Wall Street Journal URLs
func RemoveParamsFromWSJURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isWSJURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for _, param := range wsjParamsToRemove {
			if q.Has(param) {
				q.Del(param)
				changed = true
			}
		}
		if changed {
			u.RawQuery = q.Encode()
		}
		return u
	})
}
