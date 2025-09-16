package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Common tracking parameters that should be removed from all URLs
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
	if isUTMParam(param) {
		return true
	}

	for _, p := range CommonTrackingParams {
		if p == param {
			return true
		}
	}
	return false
}

// parseFragmentParams parses fragment content that contains URL-style parameters
func parseFragmentParams(fragment string) (url.Values, error) {
	// If fragment looks like it contains parameters (has = and &), parse it
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
	return values.Encode()
}

// RemoveGenericTrackingParams removes common tracking parameters from all URLs
func RemoveGenericTrackingParams(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveGenericTrackingParams: failed to read input: %w", err)
	}

	codeBlockLevel := 0
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```") {
			if codeBlockLevel == 0 {
				codeBlockLevel++
			} else {
				codeBlockLevel--
			}
		}

		if codeBlockLevel == 0 {
			rxStrict := xurls.Strict()
			lines[i] = rxStrict.ReplaceAllStringFunc(line, func(match string) string {
				u, err := url.Parse(match)
				if err != nil {
					return match
				}


				// Check all parameters and remove those that are tracking parameters
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

				// Handle fragment parameters
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
							changed = true
						}
					}
				}

				if changed {
					return u.String()
				}

				return match
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("RemoveGenericTrackingParams: failed to write output: %w", err)
	}
	return nil
}
