package links

import (
	"io"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

var WalmartTrackingParams = []string{
	"adid",
	"athbdg",
	"classType",
	"cn",
	"from",
	"gclsrc",
	"veh",
	"wmlspartner",
}

// walmartAdLabelRegex matches Walmart's wlN ad label parameters (wl0–wl12, etc.)
var walmartAdLabelRegex = regexp.MustCompile(`^wl\d+$`)

func isWalmartURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "walmart.com" || strings.HasSuffix(hostname, ".walmart.com")
}

func isWalmartTrackingParam(param string) bool {
	return isUTMParam(param) ||
		slices.Contains(WalmartTrackingParams, param) ||
		walmartAdLabelRegex.MatchString(param)
}

func RemoveParamsFromWalmartURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isWalmartURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isWalmartTrackingParam(param) {
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
