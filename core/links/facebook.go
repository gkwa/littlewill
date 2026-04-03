package links

import (
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"mvdan.cc/xurls/v2"
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
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromFacebookURLs: failed to read input: %w", err)
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

				if !isFacebookURL(u) {
					return match
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
					return u.String()
				}

				return match
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("RemoveParamsFromFacebookURLs: failed to write output: %w", err)
	}

	return nil
}
