package links

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

var textFragmentRegex = regexp.MustCompile(`(?i)^:~:text=`)

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

func isSubstackURL(u *url.URL) bool {
	return strings.HasSuffix(strings.ToLower(u.Hostname()), ".substack.com")
}

func isTheSweeklyURL(u *url.URL) bool {
	return strings.HasSuffix(strings.ToLower(u.Hostname()), "thesweekly.com")
}

var theSweeklyParamsToRemove = []string{
	"utm_source",
	"publication_id",
	"post_id",
	"utm_campaign",
	"isFreemail",
	"r",
	"triedRedirect",
	"utm_medium",
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

func RemoveParamsFromSubstackURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if isSubstackURL(u) {
			u.RawQuery = ""
		}
		return u
	})
}

func RemoveParamsFromTheSweeklyURLs(r io.Reader, w io.Writer) error {
	return processURLs(r, w, func(u *url.URL) *url.URL {
		if isTheSweeklyURL(u) {
			q := u.Query()
			for _, param := range theSweeklyParamsToRemove {
				q.Del(param)
			}
			u.RawQuery = q.Encode()
		}
		return u
	})
}

func RemoveTextFragmentsFromURLs(r io.Reader, w io.Writer) error {
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

				u = processor(u)

				return u.String()
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("processURLs: failed to write output: %w", err)
	}

	return nil
}
