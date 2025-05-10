package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Common tracking parameters that should be removed from URLs
var GenericParamsToRemove = []string{
	"_gl",
	"campaign",
	"fbclid",
	"gclid",
	"mc_cid",
	"mc_eid",
	"medium",
	"ncid",
	"ocid",
	"source",
	"srsltid",
	"utm_campaign",
	"utm_content",
	"utm_id",
	"utm_medium",
	"utm_source",
	"utm_term",
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

				q := u.Query()
				changed := false

				for _, param := range GenericParamsToRemove {
					if q.Has(param) {
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
		return fmt.Errorf("RemoveGenericTrackingParams: failed to write output: %w", err)
	}
	return nil
}
