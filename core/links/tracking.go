package links

import (
	"fmt"
	"io"
	"net/url"
	"slices"
	"sort"
	"strings"
)

// CommonTrackingParams are common tracking parameters that should be removed from all URLs
var CommonTrackingParams = []string{
	"_bhlid",
	"_ga",
	"_ga_ECJJ2Q2SJQ",
	"_gl",
	"campaign",
	"cid",
	"ck_subscriber_id",
	"fbclid",
	"gad_campaignid",
	"gad_source",
	"gbraid",
	"gclid",
	"growAuthSource",
	"growSource",
	"growUnverifiedReaderId",
	"mc_cid",
	"mc_eid",
	"medium",
	"mkt_tok",
	"ncid",
	"ocid",
	"ref",
	"scid",
	"sh_kit",
	"share_id",
	"skip_click_tracking",
	"source",
	"srsltid",
}

// isTrackingParam checks if a parameter should be removed (either UTM or in the common tracking list)
func isTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(CommonTrackingParams, param)
}

// parseFragmentParams parses fragment content that contains URL-style parameters
func parseFragmentParams(fragment string) (url.Values, error) {
	if strings.Contains(fragment, "=") {
		values, err := url.ParseQuery(fragment)
		if err != nil {
			return nil, err
		}
		return values, nil
	}
	return nil, nil
}

// buildFragmentFromParams rebuilds fragment from cleaned parameters
func buildFragmentFromParams(values url.Values) string {
	if len(values) == 0 {
		return ""
	}

	var pairs []string
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, v := range values[k] {
			if v == "" {
				pairs = append(pairs, url.QueryEscape(k))
			} else {
				pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
			}
		}
	}

	return strings.Join(pairs, "&")
}

// RemoveGenericTrackingParams removes common tracking parameters from all URLs
func RemoveGenericTrackingParams(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		q := u.Query()
		changed := false
		for param := range q {
			if isTrackingParam(param) {
				q.Del(param)
				changed = true
			}
		}
		if changed {
			u.RawQuery = q.Encode()
		}

		if u.Fragment != "" {
			fragmentParams, err := parseFragmentParams(u.Fragment)
			if err == nil && fragmentParams != nil {
				fragmentChanged := false
				for param := range fragmentParams {
					if isTrackingParam(param) {
						fragmentParams.Del(param)
						fragmentChanged = true
					}
				}
				if fragmentChanged {
					u.Fragment = buildFragmentFromParams(fragmentParams)
				}
			}
		}

		return u
	})
}
