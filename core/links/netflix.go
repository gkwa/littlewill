package links

import (
	"io"
	"net/url"
	"strings"
)

// netflixParamsToRemove are Netflix parameters that should be removed
var netflixParamsToRemove = []string{
	"s",
	"trkid",
	"shareType",
	"shareUuid",
	"trg",
	"unifiedEntityIdEncoded",
	"vlang",
	"clip",
}

// isNetflixURL checks if a URL is from Netflix
func isNetflixURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "netflix.com" || strings.HasSuffix(hostname, ".netflix.com")
}

// RemoveParamsFromNetflixURLs removes tracking parameters from Netflix URLs
func RemoveParamsFromNetflixURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isNetflixURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for _, param := range netflixParamsToRemove {
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
