package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// RedditSpecificTrackingParams are Reddit-specific params (UTM parameters are handled by shared logic)
var RedditSpecificTrackingParams = []string{
	// Branch.io parameters (both encoded and decoded versions)
	"%243p",
	"$3p",
	"%24deep_link",
	"$deep_link",
	"_branch_match_id",
	"_branch_referrer",
	// Analytics parameters
	"cId",
	"correlation_id",
	"iId",
	"post_fullname",
	"post_index",
	// Marketing parameters
	"ref_campaign",
	"ref_source",
	// Legacy Reddit parameters
	"share_id",
	"target_user",
}

// isRedditURL checks if a URL is from Reddit
func isRedditURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "reddit.com" ||
		strings.HasSuffix(hostname, ".reddit.com") ||
		hostname == "redd.it" ||
		strings.HasSuffix(hostname, ".redd.it")
}

// isRedditTrackingParam checks if a parameter should be removed from Reddit URLs
func isRedditTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(RedditSpecificTrackingParams, param)
}

// RemoveParamsFromRedditURLs removes tracking parameters from Reddit URLs
func RemoveParamsFromRedditURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isRedditURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isRedditTrackingParam(param) {
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
