package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// ShopifyTrackingParams are Shopify product recommendation tracking parameters
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
	return strings.HasSuffix(hostname, "shopify.com")
}

// isShopifyTrackingParam checks if a parameter should be removed from Shopify URLs
func isShopifyTrackingParam(param string) bool {
	return isUTMParam(param) || slices.Contains(ShopifyTrackingParams, param)
}

// RemoveParamsFromShopifyURLs removes tracking parameters from Shopify URLs
func RemoveParamsFromShopifyURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isShopifyURL(u) {
			return u
		}
		q := u.Query()
		changed := false
		for param := range q {
			if isShopifyTrackingParam(param) {
				q.Del(param)
				changed = true
			}
		}
		if changed {
			u.RawQuery = q.Encode()
		}
		return u
	})
}
