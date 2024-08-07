package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

var excludePatterns = []string{
	"google.com/maps/",
}

var excludeParams = []string{
	"q",
	"tbm",
}

var ParamsToRemove = []string{
	"bih",
	"biw",
	"client",
	"dpr",
	"ei",
	"fbs",
	"gs_lcp",
	"gs_lcrp",
	"gs_lp",
	"gs_ssp",
	"hl",
	"ictx",
	"ie",
	"oq",
	"prmd",
	"sa",
	"sca_esv",
	"sca_upv",
	"sclient",
	"source",
	"sourceid",
	"spell",
	"sqi",
	"sxsrf",
	"uact",
	// "udm", // keep image search links in tact
	"uds",
	"ved",
}

func RemoveParamsFromGoogleURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromGoogleLinks: failed to read input: %w", err)
	}

	rxStrict := xurls.Strict()
	cleaned := rxStrict.ReplaceAllFunc(buf, func(match []byte) []byte {
		urlStr := string(match)
		if isExcludedURL(urlStr) {
			return match
		}

		if strings.Contains(strings.ToLower(urlStr), "google.com") {
			cleanedURL, _, err := cleanGoogleURL(urlStr)
			if err != nil {
				return match
			}
			return []byte(cleanedURL)
		}

		return match
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromGoogleLinks: failed to write output: %w", err)
	}

	return nil
}

func isExcludedURL(urlStr string) bool {
	for _, pattern := range excludePatterns {
		if strings.Contains(strings.ToLower(urlStr), pattern) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func cleanGoogleURL(urlStr string) (string, []string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", nil, err
	}

	q := u.Query()

	var remainingParams []string

	for param := range q {
		if contains(excludeParams, param) || !contains(ParamsToRemove, param) {
			remainingParams = append(remainingParams, param)
			continue
		}
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u.String(), remainingParams, nil
}
