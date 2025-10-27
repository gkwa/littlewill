package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Mailchimp-specific tracking parameters that should be removed
// These parameters are used by Mailchimp for email campaign tracking and analytics
var MailchimpTrackingParams = []string{
	"e", // Subscriber email hash used to track engagement from email campaigns
}

// isMailchimpURL checks if a URL is from a Mailchimp domain
// Mailchimp uses multiple domains including mailchi.mp for redirects and email links
func isMailchimpURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "mailchi.mp" ||
		strings.HasSuffix(hostname, ".mailchi.mp") ||
		hostname == "mailchimp.com" ||
		strings.HasSuffix(hostname, ".mailchimp.com")
}

// isMailchimpTrackingParam checks if a parameter should be removed from Mailchimp URLs
func isMailchimpTrackingParam(param string) bool {
	// Reuse shared UTM logic first
	if isUTMParam(param) {
		return true
	}

	// Check Mailchimp-specific parameters
	for _, p := range MailchimpTrackingParams {
		if p == param {
			return true
		}
	}

	return false
}

// RemoveParamsFromMailchimpURLs removes tracking parameters from Mailchimp URLs
// Mailchimp uses the 'e' parameter to track which subscriber clicked an email link,
// along with standard UTM parameters for campaign tracking
func RemoveParamsFromMailchimpURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromMailchimpURLs: failed to read input: %w", err)
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

				if isMailchimpURL(u) {
					q := u.Query()
					changed := false

					// Use shared logic for parameter removal
					for param := range q {
						if isMailchimpTrackingParam(param) {
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
		return fmt.Errorf("RemoveParamsFromMailchimpURLs: failed to write output: %w", err)
	}

	return nil
}
