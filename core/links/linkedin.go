package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// LinkedInSpecificTrackingParams are LinkedIn-specific params (UTM parameters are handled by shared logic)
var LinkedInSpecificTrackingParams = []string{
	"rcm",
}

// isLinkedInURL checks if a URL is from LinkedIn
func isLinkedInURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "linkedin.com")
}

// isLinkedInTrackingParam checks if a parameter should be removed from LinkedIn URLs
func isLinkedInTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(LinkedInSpecificTrackingParams, param)
}

// RemoveParamsFromLinkedInURLs removes tracking parameters from LinkedIn URLs
func RemoveParamsFromLinkedInURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isLinkedInURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isLinkedInTrackingParam(param) {
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
