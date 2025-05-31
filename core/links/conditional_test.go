package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveConditionalParams(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with all conditional parameters",
			input:    "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?isFreemail=true&post_id=161256748&publication_id=2391457&r=21036&triedRedirect=true",
			expected: "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?post_id=161256748&publication_id=2391457",
		},
		{
			name:     "URL missing one conditional parameter",
			input:    "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?isFreemail=true&post_id=161256748&publication_id=2391457&r=21036",
			expected: "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?isFreemail=true&post_id=161256748&publication_id=2391457&r=21036",
		},
		{
			name:     "URL with only one conditional parameter",
			input:    "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?post_id=161256748&publication_id=2391457&r=21036",
			expected: "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?post_id=161256748&publication_id=2391457&r=21036",
		},
		{
			name:     "URL with no conditional parameters",
			input:    "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?post_id=161256748&publication_id=2391457",
			expected: "https://newsletter.systemdesignclassroom.com/p/most-systems-get-consistency-wrong?post_id=161256748&publication_id=2391457",
		},
		{
			name:     "Different domain with all conditional parameters",
			input:    "https://example.com/article?isFreemail=true&r=21036&triedRedirect=true&other=value",
			expected: "https://example.com/article?other=value",
		},
		{
			name:     "Different domain missing one conditional parameter",
			input:    "https://example.com/article?isFreemail=true&r=21036&other=value",
			expected: "https://example.com/article?isFreemail=true&r=21036&other=value",
		},
		{
			name: "Multiple URLs in text",
			input: `Check out these links:
https://newsletter.systemdesignclassroom.com/p/article1?isFreemail=true&post_id=123&r=456&triedRedirect=true
https://example.com/p/article2?isFreemail=true&post_id=789&r=101112`,
			expected: `Check out these links:
https://newsletter.systemdesignclassroom.com/p/article1?post_id=123
https://example.com/p/article2?isFreemail=true&post_id=789&r=101112`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://newsletter.systemdesignclassroom.com/p/article?isFreemail=true&r=456&triedRedirect=true

` + "```" + `
var url = "https://newsletter.systemdesignclassroom.com/p/article?isFreemail=true&r=456&triedRedirect=true";
` + "```" + `

Another URL: https://example.com/p/another?isFreemail=true&r=789&triedRedirect=true`,
			expected: `Check this URL: https://newsletter.systemdesignclassroom.com/p/article

` + "```" + `
var url = "https://newsletter.systemdesignclassroom.com/p/article?isFreemail=true&r=456&triedRedirect=true";
` + "```" + `

Another URL: https://example.com/p/another`,
		},
		{
			name:     "URL with conditional parameters in different order",
			input:    "https://another-site.com/page?triedRedirect=true&other=keep&isFreemail=false&r=12345",
			expected: "https://another-site.com/page?other=keep",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveConditionalParams(input, &output)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			result := output.String()
			if diff := cmp.Diff(tc.expected, result); diff != "" {
				t.Errorf("Unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
