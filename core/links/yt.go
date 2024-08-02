package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

type URLProcessor interface {
	Process(url *url.URL) *url.URL
}

type YouTubeURLProcessor struct{}

func (p *YouTubeURLProcessor) Process(u *url.URL) *url.URL {
	if isYouTubeURL(u) {
		q := u.Query()
		q.Del("si")
		q.Del("app")
		u.RawQuery = q.Encode()
	}
	return u
}

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

func RemoveYoutubeParams(r io.Reader, w io.Writer) error {
	processors := []URLProcessor{
		&YouTubeURLProcessor{},
	}
	return processURLs(r, w, processors...)
}

func processURLs(r io.Reader, w io.Writer, processors ...URLProcessor) error {
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

		for _, processor := range processors {
			u = processor.Process(u)
		}

		return []byte(u.String())
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("processURLs: failed to write output: %w", err)
	}

	return nil
}
