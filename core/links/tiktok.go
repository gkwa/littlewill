package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// TikTokSpecificTrackingParams are TikTok-specific params kept separate from the
// generic list because "t" is generic enough to be legitimate on other sites
// (e.g. timestamps for YouTube start positions).
var TikTokSpecificTrackingParams = []string{
	"t",
}

func isTikTokURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "tiktok.com" || strings.HasSuffix(hostname, ".tiktok.com")
}

func isTikTokTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(TikTokSpecificTrackingParams, param)
}

func RemoveParamsFromTikTokURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isTikTokURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isTikTokTrackingParam(param) {
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
