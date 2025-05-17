package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// TechCrunch email parameters that should be removed
var techcrunchParamsToRemove = []string{
	"ecid",
	"utm_campaign",
	"utm_medium",
	"_hsenc",
	"_hsmi",
	"utm_content",
	"utm_source",
}

// isTechCrunchURL checks if a URL is from TechCrunch's email domain
func isTechCrunchURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "techcrunch.com")
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

					for _, param := range techcrunchParamsToRemove {
						if q.Has(param) {
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
