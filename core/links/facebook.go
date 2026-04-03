package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// FacebookSpecificTrackingParams are Facebook-specific params kept separate from the
// generic list because names like "referral_code" and "tracking" are generic enough
// to be legitimate on other sites (e.g. referral_code for discounts, tracking for shipments).
var FacebookSpecificTrackingParams = []string{
	"in_reels_tab_context",
	"referral_code",
	"referral_source",
	"referral_story_type",
	"surface_type",
	"tracking",
}

// isFacebookURL checks if a URL is from Facebook
func isFacebookURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "facebook.com")
}

// isFacebookTrackingParam checks if a parameter should be removed from Facebook URLs
func isFacebookTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(FacebookSpecificTrackingParams, param)
}

// RemoveParamsFromFacebookURLs removes tracking parameters from Facebook URLs
func RemoveParamsFromFacebookURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isFacebookURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isFacebookTrackingParam(param) {
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
