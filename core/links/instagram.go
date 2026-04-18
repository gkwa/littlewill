package links

import (
	"io"
	"net/url"
	"strings"
)

var instagramParamsToRemove = []string{
	"hl",
	"igsh",
	"igshid",
}

func isInstagramURL(u *url.URL) bool {
	hostname := strings.ToLower(u.Hostname())
	return hostname == "instagram.com" || strings.HasSuffix(hostname, ".instagram.com")
}

func RemoveParamsFromInstagramURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isInstagramURL(u) {
			return u
		}
		q := u.Query()
		for _, param := range instagramParamsToRemove {
			q.Del(param)
		}
		if len(q) == 0 {
			u.RawQuery = ""
		} else {
			u.RawQuery = q.Encode()
		}
		return u
	})
}
