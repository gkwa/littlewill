package links

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	markdownLinkRegex          = regexp.MustCompile(`\[([^][]+)\](\(((?:[^()]+|\([^()]*\))+)\))`)
	markdownLinkWithTitleRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\s+"[^"]+"\)`)
	whitespaceRegex            = regexp.MustCompile(`[\s\t\n\r]+`)
)

func RemoveWhitespaceFromMarkdownLinks(r io.Reader, w io.Writer) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("CleanupMarkdownLinks: failed to read input: %w", err)
	}

	cleaned := markdownLinkRegex.ReplaceAllFunc(buf, func(match []byte) []byte {
		submatches := markdownLinkRegex.FindSubmatch(match)
		if len(submatches) >= 4 {
			description := whitespaceRegex.ReplaceAllString(string(submatches[1]), " ")
			description = strings.TrimSpace(description)
			url := string(submatches[3])
			return []byte(fmt.Sprintf("[%s](%s)", description, url))
		}
		return match
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
