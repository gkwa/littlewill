package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// WSJ parameters that should be removed
var wsjParamsToRemove = []string{
	"mod",
	"reflink",
	"ref",
	"st",
}

// isWSJURL checks if a URL is from Wall Street Journal
func isWSJURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "wsj.com")
}

// RemoveParamsFromWSJURLs removes tracking parameters from Wall Street Journal URLs
func RemoveParamsFromWSJURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromWSJURLs: failed to read input: %w", err)
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

				if isWSJURL(u) {
					q := u.Query()
					changed := false

					for _, param := range wsjParamsToRemove {
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
		return fmt.Errorf("RemoveParamsFromWSJURLs: failed to write output: %w", err)
	}

	return nil
}
