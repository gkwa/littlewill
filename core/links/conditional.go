package links

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// ConditionalParamGroup represents a group of parameters that should only be removed
// if all parameters in the group are present in the URL
type ConditionalParamGroup struct {
	Params []string // Parameters that must all be present to be removed
}

// RemoveConditionalParams removes parameters only when all parameters in a group are present
func RemoveConditionalParams(r io.Reader, w io.Writer) error {
	// Define conditional parameter groups
	conditionalGroups := []ConditionalParamGroup{
		{
			Params: []string{"isFreemail", "r", "triedRedirect", "post_id", "publication_id"},
		},
		// Add more conditional groups here as needed
		// {
		//     Params: []string{"param1", "param2", "param3"},
		// },
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("RemoveConditionalParams: failed to read input: %w", err)
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

				// Check each conditional group
				for _, group := range conditionalGroups {
					if shouldRemoveParams(u, group.Params) {
						q := u.Query()
						for _, param := range group.Params {
							q.Del(param)
						}
						u.RawQuery = q.Encode()
					}
				}

				return u.String()
			})
		}
	}

	_, err = w.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return fmt.Errorf("RemoveConditionalParams: failed to write output: %w", err)
	}
	return nil
}

// shouldRemoveParams checks if all parameters in the group are present in the URL
func shouldRemoveParams(u *url.URL, params []string) bool {
	q := u.Query()
	for _, param := range params {
		if !q.Has(param) {
			return false
		}
	}
	return true
}
