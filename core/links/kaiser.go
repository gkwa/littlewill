package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Kaiser Permanente-specific tracking parameters that should be removed
var KaiserSpecificTrackingParams = []string{
	"promo_id",
	"wt.mc_id",
	"wt.tsrc",
	"cid",
	"ad_id",
	"cat",
	"mkt_tok",
}

// isKaiserURL checks if a URL is from Kaiser Permanente
func isKaiserURL(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Hostname()), "kaiserpermanente.org")
}

// isKaiserTrackingParam checks if a parameter should be removed from Kaiser URLs
func isKaiserTrackingParam(param string) bool {
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}

	// Check Kaiser-specific parameters
	for _, p := range KaiserSpecificTrackingParams {
		if p == param {
			return true
		}
	}
	return false
}

// RemoveParamsFromKaiserURLs removes tracking parameters from Kaiser Permanente URLs
func RemoveParamsFromKaiserURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromKaiserURLs: failed to read input: %w", err)
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

				if isKaiserURL(u) {
					q := u.Query()
					changed := false

					// Use shared logic for parameter removal
					for param := range q {
						if isKaiserTrackingParam(param) {
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
		return fmt.Errorf("RemoveParamsFromKaiserURLs: failed to write output: %w", err)
	}

	return nil
}
