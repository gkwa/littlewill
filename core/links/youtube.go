package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

func RemoveParamsFromYouTubeURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromYouTubeLinks: failed to read input: %w", err)
	}

	youTubeDomains := []string{
		"youtube.com",
		"youtu.be",
	}

	rxStrict := xurls.Strict()
	cleaned := rxStrict.ReplaceAllFunc(buf, func(match []byte) []byte {
		u, err := url.Parse(string(match))
		if err != nil {
			return match
		}

		for _, domain := range youTubeDomains {
			if strings.Contains(strings.ToLower(u.Hostname()), domain) {
				q := u.Query()
				q.Del("si")
				u.RawQuery = q.Encode()
				return []byte(u.String())
			}
		}

		return match
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromYouTubeLinks: failed to write output: %w", err)
	}

	return nil
}
