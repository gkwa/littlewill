package links

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

var textFragmentRegex = regexp.MustCompile(`^:~:text=`)

func isYouTubeURL(u *url.URL) bool {
	youTubeDomains := []string{
		"youtube.com",
		"youtu.be",
	}
	for _, domain := range youTubeDomains {
		if strings.Contains(strings.ToLower(u.Hostname()), domain) {
			return true
		}
	}
	return false
}

func RemoveParamsFromYouTubeURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if isYouTubeURL(u) {
			q := u.Query()
			q.Del("si")
			q.Del("app")
			u.RawQuery = q.Encode()
		}
		return u
	})
}

func RemoveTextFragments(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if isTextFragment(u.Fragment) {
			u.Fragment = ""
		}
		return u
	})
}

func isTextFragment(fragment string) bool {
	return textFragmentRegex.MatchString(fragment)
}

func processURLs(r io.Reader, w io.Writer, processor func(*url.URL) *url.URL) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("processURLs: failed to read input: %w", err)
	}

	rxStrict := xurls.Strict()
	cleaned := rxStrict.ReplaceAllFunc(buf, func(match []byte) []byte {
		u, err := url.Parse(string(match))
		if err != nil {
			return match
		}

		u = processor(u)

		return []byte(u.String())
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("processURLs: failed to write output: %w", err)
	}

	return nil
}
