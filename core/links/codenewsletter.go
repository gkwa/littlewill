package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// CodeNewsletterSpecificTrackingParams are codenewsletter.ai-specific params
// (UTM parameters are handled by shared logic)
var CodeNewsletterSpecificTrackingParams = []string{
	"_bhlid",
	"jwt_token",
}

// isCodeNewsletterURL checks if a URL is from codenewsletter.ai
func isCodeNewsletterURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "codenewsletter.ai" || strings.HasSuffix(hostname, ".codenewsletter.ai")
}

// isCodeNewsletterTrackingParam checks if a parameter should be removed from codenewsletter.ai URLs
func isCodeNewsletterTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(CodeNewsletterSpecificTrackingParams, param)
}

// RemoveParamsFromCodeNewsletterURLs removes tracking parameters from codenewsletter.ai URLs
func RemoveParamsFromCodeNewsletterURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isCodeNewsletterURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isCodeNewsletterTrackingParam(param) {
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
