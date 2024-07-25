package core

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

var (
	markdownLinkRegex          = regexp.MustCompile(`\[\s*(\S.*?)\s*\]\(`)
	markdownLinkWithTitleRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\s+"[^"]+"\)`)
)

func RemoveWhitespaceFromMarkdownLinks(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("CleanupMarkdownLinks: failed to read input: %w", err)
	}

	cleaned := markdownLinkRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		return markdownLinkRegex.ReplaceAll(match, []byte("[$1]("))
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("CleanupMarkdownLinks: failed to write output: %w", err)
	}

	return nil
}

func RemoveTitlesFromMarkdownLinks(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveTitlesFromMarkdownLinks: failed to read input: %w", err)
	}

	cleaned := markdownLinkWithTitleRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		return markdownLinkWithTitleRegex.ReplaceAll(match, []byte("[$1]($2)"))
	})

	_, err = w.Write(cleaned)
	if err != nil {
		return fmt.Errorf("RemoveTitlesFromMarkdownLinks: failed to write output: %w", err)
	}

	return nil
}

func RemoveParamsFromYouTubeLinks(r io.Reader, w io.Writer) error {
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

func ApplyTransforms(content []byte, transforms ...func(io.Reader, io.Writer) error) ([]byte, error) {
	var processedContent bytes.Buffer
	currentContent := bytes.NewReader(content)

	for _, transform := range transforms {
		processedContent.Reset()
		err := transform(currentContent, &processedContent)
		if err != nil {
			return nil, err
		}
		currentContent = bytes.NewReader(processedContent.Bytes())
	}

	return processedContent.Bytes(), nil
}
