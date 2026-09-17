package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// MediumSpecificTrackingParams are Medium-specific params (UTM parameters are handled by shared logic)
var MediumSpecificTrackingParams = []string{
	"source",
}

// isMediumURL checks if a URL is from medium.com or one of its subdomains
func isMediumURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "medium.com" || strings.HasSuffix(hostname, ".medium.com")
}

// RemoveParamsFromMediumURLs removes tracking parameters from Medium URLs
func RemoveParamsFromMediumURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isMediumURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if slices.Contains(MediumSpecificTrackingParams, param) {
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
