package links

import (
	"fmt"
	"io"
	"regexp"
)

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
