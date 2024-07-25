package core

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
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
