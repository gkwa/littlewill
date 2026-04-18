package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// TechCrunchSpecificTrackingParams are TechCrunch-specific params (UTM parameters are handled by shared logic)
var TechCrunchSpecificTrackingParams = []string{
	"ecid",
	"_hsenc",
	"_hsmi",
}

// isTechCrunchURL checks if a URL is from TechCrunch
func isTechCrunchURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "techcrunch.com")
}

// isTechCrunchTrackingParam checks if a parameter should be removed from TechCrunch URLs
func isTechCrunchTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(TechCrunchSpecificTrackingParams, param)
}

// RemoveParamsFromTechCrunchURLs removes tracking parameters from TechCrunch URLs
func RemoveParamsFromTechCrunchURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isTechCrunchURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isTechCrunchTrackingParam(param) {
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
