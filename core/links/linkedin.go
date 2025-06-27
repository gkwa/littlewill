package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// LinkedIn-specific tracking parameters that should be removed (UTM parameters are handled by shared logic)
var LinkedInSpecificTrackingParams = []string{
	"rcm",
}

// isLinkedInURL checks if a URL is from LinkedIn
func isLinkedInURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "linkedin.com")
}

// isLinkedInTrackingParam checks if a parameter should be removed from LinkedIn URLs
func isLinkedInTrackingParam(param string) bool {
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}

	// Check LinkedIn-specific parameters
	for _, p := range LinkedInSpecificTrackingParams {
		if p == param {
			return true
		}
	}
	return false
}

// RemoveParamsFromLinkedInURLs removes tracking parameters from LinkedIn URLs
func RemoveParamsFromLinkedInURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromLinkedInURLs: failed to read input: %w", err)
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

				if isLinkedInURL(u) {
					q := u.Query()
					changed := false

					// Use shared logic for parameter removal
					for param := range q {
						if isLinkedInTrackingParam(param) {
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
		return fmt.Errorf("RemoveParamsFromLinkedInURLs: failed to write output: %w", err)
	}

	return nil
}
