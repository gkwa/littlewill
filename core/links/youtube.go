package links

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
)

// YouTubeParamsToRemove are YouTube parameters that should be removed
var YouTubeParamsToRemove = []string{
	"si",
	"app",
	"feature",
	"sqp",
	"rs",
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
		"ytimg.com",
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
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if !isYouTubeURL(u) {
			return u
		}

		q := u.Query()
		changed := false

		for _, param := range YouTubeParamsToRemove {
			if q.Has(param) {
				q.Del(param)
				changed = true
			}
		}

		// Convert youtube.com/watch?v=X to youtu.be/X
		if strings.Contains(strings.ToLower(u.Hostname()), "youtube.com") && u.Path == "/watch" && q.Has("v") {
			videoID := q.Get("v")
			q.Del("v")
			u.Host = "youtu.be"
			u.Path = "/" + videoID
			changed = true
		}

		// Convert youtube.com/shorts/VIDEO_ID to youtu.be/VIDEO_ID
		if strings.Contains(strings.ToLower(u.Hostname()), "youtube.com") && strings.HasPrefix(u.Path, "/shorts/") {
			videoID := strings.TrimPrefix(u.Path, "/shorts/")
			u.Host = "youtu.be"
			u.Path = "/" + videoID
			changed = true
		}

		// Convert youtube.com/live/VIDEO_ID to youtu.be/VIDEO_ID
		if strings.Contains(strings.ToLower(u.Hostname()), "youtube.com") && strings.HasPrefix(u.Path, "/live/") {
			videoID := strings.TrimPrefix(u.Path, "/live/")
			u.Host = "youtu.be"
			u.Path = "/" + videoID
			changed = true
		}

		if changed {
			u.RawQuery = q.Encode()
		}
		return u
	})
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
