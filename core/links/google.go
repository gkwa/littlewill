package links

import (
	"io"
	"net"
	"net/url"
	"slices"
	"strings"
)

var excludePatterns = []string{
	"google.com/maps/",
}

var KeepParams = []string{
	"q",
	"tbm",
}

var ParamsToRemove = []string{
	"aep",
	"bih",
	"biw",
	"client",
	"cshid",
	"csuir",
	"dpr",
	"ei",
	"fbs",
	"gs_lcp",
	"gs_lcrp",
	"gs_lp",
	"gs_ssp",
	"hl",
	"hs",
	"ictx",
	"ie",
	"mstk",
	"ntc",
	"num",
	"oq",
	"prmd",
	"sa",
	"sca_esv",
	"sca_upv",
	"sclient",
	"sei",
	"si",
	"source",
	"sourceid",
	"spell",
	"sqi",
	"stick",
	"sxsrf",
	"uact",
	"uds",
	"ved",
}

func RemoveParamsFromGoogleURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if isExcludedURL(u.String()) {
			return u
		}
		if !strings.Contains(strings.ToLower(u.Hostname()), "google.com") {
			return u
		}
		cleaned, _, err := cleanGoogleURL(u.String())
		if err != nil {
			return u
		}
		parsed, err := url.Parse(cleaned)
		if err != nil {
			return u
		}
		stripWWWHost(parsed)
		return parsed
	})
}

func isExcludedURL(urlStr string) bool {
	for _, pattern := range excludePatterns {
		if strings.Contains(strings.ToLower(urlStr), pattern) {
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
		if slices.Contains(KeepParams, param) || !slices.Contains(ParamsToRemove, param) {
			remainingParams = append(remainingParams, param)
			continue
		}
		q.Del(param)
	}

	u.RawQuery = q.Encode()
	return u.String(), remainingParams, nil
}

// stripWWWHost drops the www. prefix from a Google host. Google 301-redirects
// the bare host to www., so the shorter form resolves identically.
func stripWWWHost(u *url.URL) {
	host := u.Hostname()
	if !strings.HasPrefix(strings.ToLower(host), "www.") {
		return
	}

	bare := host[len("www."):]
	if port := u.Port(); port != "" {
		u.Host = net.JoinHostPort(bare, port)
		return
	}

	u.Host = bare
}
