package links

import (
	"io"
	"net/url"
	"slices"
	"strings"
)

// AmazonTrackingParams are Amazon-specific tracking parameters that should be removed
var AmazonTrackingParams = []string{
	"_encoding",
	"content-id",
	"crid",
	"cv_ct_cx",
	"dib",
	"dib_tag",
	"ds",
	"keywords",
	"pd_rd_i",
	"pd_rd_r",
	"pd_rd_w",
	"pd_rd_wg",
	"pf_rd_p",
	"pf_rd_r",
	"psc",
	"qid",
	"ref",
	"ref_",
	"sbo",
	"smid",
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
	return isUTMParam(param) || slices.Contains(AmazonTrackingParams, param)
}

// RemoveParamsFromAmazonURLs removes tracking parameters from Amazon URLs
func RemoveParamsFromAmazonURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isAmazonURL(u) {
			return u
		}

		// Remove path segments that start with "ref="
		pathSegments := strings.Split(u.Path, "/")
		var cleanedSegments []string
		for _, segment := range pathSegments {
			if !strings.HasPrefix(segment, "ref=") {
				cleanedSegments = append(cleanedSegments, segment)
			}
		}
		u.Path = stripTrailingSlash(strings.Join(cleanedSegments, "/"))

		q := u.Query()
		for param := range q {
			if isAmazonTrackingParam(param) {
				q.Del(param)
			}
		}
		// Always re-encode to normalize parameter order
		u.RawQuery = q.Encode()
		return u
	})
}
