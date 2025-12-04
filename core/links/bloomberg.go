package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Bloomberg parameters that should be removed
var bloombergParamsToRemove = []string{
	"accessToken",
	"leadSource",
}

// isBloombergURL checks if a URL is from Bloomberg
func isBloombergURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "bloomberg.com")
}

// RemoveParamsFromBloombergURLs removes tracking parameters from Bloomberg URLs
func RemoveParamsFromBloombergURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromBloombergURLs: failed to read input: %w", err)
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

				if isBloombergURL(u) {
					q := u.Query()
					changed := false
					for _, param := range bloombergParamsToRemove {
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
		return fmt.Errorf("RemoveParamsFromBloombergURLs: failed to write output: %w", err)
	}

	return nil
}
