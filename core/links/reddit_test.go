package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromRedditURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Reddit post with share_id parameter",
			input:    "https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/?share_id=mHY19WpSs_UfZsKQhbBrx",
			expected: "https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/",
		},
		{
			name:     "Reddit short URL with share_id",
			input:    "https://redd.it/1lb7mga?share_id=mHY19WpSs_UfZsKQhbBrx",
			expected: "https://redd.it/1lb7mga",
		},
		{
			name:     "Reddit URL with UTM parameters using shared logic",
			input:    "https://www.reddit.com/r/test/comments/789?utm_source=newsletter&utm_medium=email&sort=new",
			expected: "https://www.reddit.com/r/test/comments/789?sort=new",
		},
		{
			name:     "Reddit URL with share_id and UTM parameters",
			input:    "https://www.reddit.com/r/programming/comments/123456?share_id=abc123&utm_source=newsletter&utm_medium=email",
			expected: "https://www.reddit.com/r/programming/comments/123456",
		},
		{
			name:     "Reddit URL with share_id and other parameters",
			input:    "https://www.reddit.com/r/test/comments/789?share_id=xyz&sort=new&other=keep",
			expected: "https://www.reddit.com/r/test/comments/789?other=keep&sort=new",
		},
		{
			name:     "Reddit URL without tracking parameters",
			input:    "https://www.reddit.com/r/golang/comments/456789/some_post/",
			expected: "https://www.reddit.com/r/golang/comments/456789/some_post/",
		},
		{
			name:     "Non-Reddit URL with share_id",
			input:    "https://example.com?share_id=shouldnotberemoved",
			expected: "https://example.com?share_id=shouldnotberemoved",
		},
		{
			name: "Multiple Reddit URLs in text",
			input: `Check out these Reddit posts:
https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/?share_id=mHY19WpSs_UfZsKQhbBrx
https://redd.it/456789?share_id=another123&utm_source=app`,
			expected: `Check out these Reddit posts:
https://www.reddit.com/r/Sourdough/comments/1lb7mga/crumb_read_customer_complains_my_sourdough_rises/
https://redd.it/456789`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Reddit URL: https://www.reddit.com/r/test/comments/123?share_id=abc123

` + "```" + `
var redditUrl = "https://www.reddit.com/r/test/comments/123?share_id=abc123";
` + "```" + `

Another Reddit URL: https://redd.it/456?share_id=def456&utm_source=test`,
			expected: `Check this Reddit URL: https://www.reddit.com/r/test/comments/123

` + "```" + `
var redditUrl = "https://www.reddit.com/r/test/comments/123?share_id=abc123";
` + "```" + `

Another Reddit URL: https://redd.it/456`,
		},
		{
			name:     "Reddit subdomain URL",
			input:    "https://old.reddit.com/r/programming/comments/123?share_id=test123",
			expected: "https://old.reddit.com/r/programming/comments/123",
		},
		{
			name:     "Reddit URL with custom UTM parameter",
			input:    "https://www.reddit.com/r/golang/?utm_custom_param=value&sort=new",
			expected: "https://www.reddit.com/r/golang/?sort=new",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromRedditURLs(input, &output)
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
