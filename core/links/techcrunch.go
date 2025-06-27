package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// TechCrunch-specific tracking parameters that should be removed (UTM parameters are handled by shared logic)
var TechCrunchSpecificTrackingParams = []string{
	"ecid",
	"_hsenc",
	"_hsmi",
}

// isTechCrunchURL checks if a URL is from TechCrunch's email domain
func isTechCrunchURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "techcrunch.com")
}

// isTechCrunchTrackingParam checks if a parameter should be removed from TechCrunch URLs
func isTechCrunchTrackingParam(param string) bool {
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}

	// Check TechCrunch-specific parameters
	for _, p := range TechCrunchSpecificTrackingParams {
		if p == param {
			return true
		}
	}
	return false
}

// RemoveParamsFromTechCrunchURLs removes tracking parameters from TechCrunch email URLs
func RemoveParamsFromTechCrunchURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromTechCrunchURLs: failed to read input: %w", err)
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

				if isTechCrunchURL(u) {
					q := u.Query()
					changed := false

					// Use shared logic for parameter removal
					for param := range q {
						if isTechCrunchTrackingParam(param) {
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
		return fmt.Errorf("RemoveParamsFromTechCrunchURLs: failed to write output: %w", err)
	}

	return nil
}
