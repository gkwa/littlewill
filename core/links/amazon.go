package links

import (
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Amazon-specific tracking parameters that should be removed
var AmazonTrackingParams = []string{
	"_encoding",
	"content-id",
	"crid",
	"cv_ct_cx",
	"dib",
	"dib_tag",
	"keywords",
	"pd_rd_i",
	"pd_rd_r",
	"pd_rd_w",
	"pd_rd_wg",
	"pf_rd_p",
	"pf_rd_r",
	"qid",
	"ref",
	"ref_",
	"sbo",
	"sp_csd",
	"sprefix",
	"sr",
	"th",
}

// isAmazonURL checks if a URL is from Amazon
func isAmazonURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "amazon.com" ||
		strings.HasSuffix(hostname, ".amazon.com") ||
		strings.Contains(hostname, "amazon.") || // Handles amazon.co.uk, amazon.de, etc.
		hostname == "amzn.to" ||
		strings.HasSuffix(hostname, ".amzn.to")
}

// isAmazonTrackingParam checks if a parameter should be removed from Amazon URLs
func isAmazonTrackingParam(param string) bool {
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}

	// Check Amazon-specific parameters
	return slices.Contains(AmazonTrackingParams, param)
}

// RemoveParamsFromAmazonURLs removes tracking parameters from Amazon URLs
func RemoveParamsFromAmazonURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromAmazonURLs: failed to read input: %w", err)
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

				if isAmazonURL(u) {
					// Remove path segments that start with "ref="
					pathSegments := strings.Split(u.Path, "/")
					var cleanedSegments []string
					for _, segment := range pathSegments {
						if !strings.HasPrefix(segment, "ref=") {
							cleanedSegments = append(cleanedSegments, segment)
						}
					}
					u.Path = strings.Join(cleanedSegments, "/")

					q := u.Query()

					// Use shared logic for parameter removal
					for param := range q {
						if isAmazonTrackingParam(param) {
							q.Del(param)
						}
					}

					// Always re-encode to normalize parameter order
					u.RawQuery = q.Encode()
					return u.String()
				}

				return match
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("RemoveParamsFromAmazonURLs: failed to write output: %w", err)
	}

	return nil
}
