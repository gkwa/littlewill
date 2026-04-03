package links

import (
	"io"
	"net/url"
)

// ConditionalParamGroup represents a group of parameters that should only be removed
// if all parameters in the group are present in the URL
type ConditionalParamGroup struct {
	Params []string // Parameters that must all be present to be removed
}

// RemoveConditionalParams removes parameters only when all parameters in a group are present
func RemoveConditionalParams(r io.Reader, w io.Writer) error {
	// These parameters must all be present to be removed
	conditionalGroups := []ConditionalParamGroup{
		{
			Params: []string{"isFreemail", "r", "triedRedirect"},
		},
		// Add more conditional groups here as needed
		// {
		//     Params: []string{"param1", "param2", "param3"},
		// },
	}

	return processURLs(r, w, func(u *url.URL) *url.URL {
		for _, group := range conditionalGroups {
			if shouldRemoveParams(u, group.Params) {
				q := u.Query()
				for _, param := range group.Params {
					q.Del(param)
				}
				u.RawQuery = q.Encode()
			}
		}
		return u
	})
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
