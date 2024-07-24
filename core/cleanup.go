package core

import (
	"fmt"
	"io"
	"regexp"
)

var markdownLinkRegex = regexp.MustCompile(`\[\s*(\S.*?)\s*\]\(`)

func CleanupMarkdownLinks(r io.Reader, w io.Writer) error {
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
