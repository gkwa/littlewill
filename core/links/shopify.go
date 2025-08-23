package links

import (
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"mvdan.cc/xurls/v2"
)

// Shopify product recommendation tracking parameters that should be removed
var ShopifyTrackingParams = []string{
	"pr_prod_strat", // Product recommendation strategy tracking parameter
	"pr_rec_id",     // Product recommendation ID for tracking purposes
	"pr_rec_pid",    // Product recommendation product ID tracking
	"pr_ref_pid",    // Product recommendation reference product ID
	"pr_seq",        // Product recommendation sequence tracking
}

// isShopifyURL checks if a URL is from a Shopify store
func isShopifyURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	// Check for myshopify.com domains and other common Shopify patterns
	return strings.HasSuffix(hostname, "shopify.com")
}

// isShopifyTrackingParam checks if a parameter should be removed from Shopify URLs
func isShopifyTrackingParam(param string) bool {
	// Reuse shared UTM logic
	if isUTMParam(param) {
		return true
	}
	// Check Shopify-specific parameters
	return slices.Contains(ShopifyTrackingParams, param)
}

// RemoveParamsFromShopifyURLs removes tracking parameters from Shopify URLs
func RemoveParamsFromShopifyURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromShopifyURLs: failed to read input: %w", err)
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

				if isShopifyURL(u) {
					q := u.Query()
					changed := false
					// Use shared logic for parameter removal
					for param := range q {
						if isShopifyTrackingParam(param) {
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
		return fmt.Errorf("RemoveParamsFromShopifyURLs: failed to write output: %w", err)
	}

	return nil
}
