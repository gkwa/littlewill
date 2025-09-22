package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Reddit-specific tracking parameters that should be removed (UTM parameters are handled by shared logic)
var RedditSpecificTrackingParams = []string{
	// Branch.io parameters (both encoded and decoded versions)
	"%243p",
	"$3p",
	"%24deep_link",
	"$deep_link",
	"_branch_match_id",
	"_branch_referrer",
	// Analytics parameters
	"correlation_id",
	"post_fullname",
	"post_index",
	// Marketing parameters
	"ref_campaign",
	"ref_source",
	// Legacy Reddit parameters
	"share_id",
	"cId",
	"iId",
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
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}

	// Check Reddit-specific parameters
	for _, p := range RedditSpecificTrackingParams {
		if p == param {
			return true
		}
	}

	return false
}

// RemoveParamsFromRedditURLs removes tracking parameters from Reddit URLs
// Uses shared UTM detection logic plus Reddit-specific parameters
func RemoveParamsFromRedditURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromRedditURLs: failed to read input: %w", err)
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

				if isRedditURL(u) {
					q := u.Query()
					changed := false

					// Use shared logic for parameter removal
					for param := range q {
						if isRedditTrackingParam(param) {
							q.Del(param)
							changed = true
						}
					}

					if changed {
						u.RawQuery = q.Encode()
						return u.String()
					}
				}

				return match
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("RemoveParamsFromRedditURLs: failed to write output: %w", err)
	}

	return nil
}
