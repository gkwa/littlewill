package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Netflix parameters that should be removed
var netflixParamsToRemove = []string{
	"s",
	"trkid",
	"shareType",
	"shareUuid",
	"trg",
	"unifiedEntityIdEncoded",
	"vlang",
	"clip",
}

// isNetflixURL checks if a URL is from Netflix
func isNetflixURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "netflix.com" || strings.HasSuffix(hostname, ".netflix.com")
}

// RemoveParamsFromNetflixURLs removes tracking parameters from Netflix URLs
func RemoveParamsFromNetflixURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromNetflixURLs: failed to read input: %w", err)
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

				if isNetflixURL(u) {
					q := u.Query()
					changed := false

					for _, param := range netflixParamsToRemove {
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
		return fmt.Errorf("RemoveParamsFromNetflixURLs: failed to write output: %w", err)
	}

	return nil
}
