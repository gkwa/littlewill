package links

import (
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

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
}
