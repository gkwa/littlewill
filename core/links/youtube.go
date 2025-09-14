package links

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

// YouTube parameters that should be removed
var YouTubeParamsToRemove = []string{
	"si",
	"app",
}

var youtubeCountRegex *regexp.Regexp

func init() {
	var err error
	youtubeCountRegex, err = buildYouTubeCountRegex()
	if err != nil {
		panic(fmt.Sprintf("Failed to build YouTube count regex: %v", err))
	}
}

func buildYouTubeCountRegex() (*regexp.Regexp, error) {
	count := `\(\d+\)`
	title := `.+?`
	youtubeURLs := `(?:youtube\.com|youtu\.be)`
	url := `https?://(?:www\.)?` + youtubeURLs + `/[^\s]+`
	optionalTitle := `(?:\s+"[^"]*")?`

	pattern := fmt.Sprintf(`\[\s*(%s)\s*(%s)\]\((%s)(%s)\)`, count, title, url, optionalTitle)

	return regexp.Compile(pattern)
}

// isYouTubeURL checks if a URL is from YouTube
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

// RemoveParamsFromYouTubeURLs removes tracking parameters from YouTube URLs
func RemoveParamsFromYouTubeURLs(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveParamsFromYouTubeURLs: failed to read input: %w", err)
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

				if isYouTubeURL(u) {
					q := u.Query()
					changed := false

					for _, param := range YouTubeParamsToRemove {
						if q.Has(param) {
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
		return fmt.Errorf("RemoveParamsFromYouTubeURLs: failed to write output: %w", err)
	}

	return nil
}

// RemoveYouTubeCountFromMarkdownLinks removes view counts from YouTube markdown links
func RemoveYouTubeCountFromMarkdownLinks(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveYouTubeCountFromLinks: failed to read input: %w", err)
	}

	cleaned := youtubeCountRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		return youtubeCountRegex.ReplaceAll(match, []byte("[$2]($3$4)"))
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("RemoveYouTubeCountFromLinks: failed to write output: %w", err)
	}

	return nil
}
