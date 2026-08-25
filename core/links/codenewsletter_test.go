package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromCodeNewsletterURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Newsletter post with utm, _bhlid and jwt_token",
			input:    "https://codenewsletter.ai/p/agentsky-ships-one-api-for-every-coding-agent?utm_source=codenewsletter.ai&utm_medium=newsletter&utm_campaign=agentsky-ships-one-api-for-every-coding-agent&_bhlid=05cd75c8a51f813b55d4bb9b0d06e67caba4519b&jwt_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWJzY3JpYmVyX2lkIjoiZTM3ZmY2MTgtMGVkMS00MTI0LTgyMTctN2FkMjA0MDkwYWRjIn0.LXi8tA92X6quCs8Q5sA73mGjGzZq2u9pXGxD7ftBd0w",
			expected: "https://codenewsletter.ai/p/agentsky-ships-one-api-for-every-coding-agent",
		},
		{
			name:     "Newsletter post with only _bhlid",
			input:    "https://codenewsletter.ai/p/some-post?_bhlid=abc123",
			expected: "https://codenewsletter.ai/p/some-post",
		},
		{
			name:     "Newsletter post keeps unknown parameters",
			input:    "https://codenewsletter.ai/p/some-post?_bhlid=abc123&page=2",
			expected: "https://codenewsletter.ai/p/some-post?page=2",
		},
		{
			name:     "Newsletter post without tracking parameters",
			input:    "https://codenewsletter.ai/p/some-post",
			expected: "https://codenewsletter.ai/p/some-post",
		},
		{
			name:     "Trailing slash gets stripped",
			input:    "https://codenewsletter.ai/p/some-post/",
			expected: "https://codenewsletter.ai/p/some-post",
		},
		{
			name:     "Subdomain URL",
			input:    "https://www.codenewsletter.ai/p/some-post?utm_source=twitter&id=7",
			expected: "https://www.codenewsletter.ai/p/some-post?id=7",
		},
		{
			name:     "Non-codenewsletter URL with similar parameters",
			input:    "https://example.com/article?_bhlid=abc123&jwt_token=xyz",
			expected: "https://example.com/article?_bhlid=abc123&jwt_token=xyz",
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Read this: https://codenewsletter.ai/p/post-one?_bhlid=abc123
` + "```" + `
var url = "https://codenewsletter.ai/p/post-one?_bhlid=abc123";
` + "```" + `
And this: https://codenewsletter.ai/p/post-two?utm_medium=newsletter`,
			expected: `Read this: https://codenewsletter.ai/p/post-one
` + "```" + `
var url = "https://codenewsletter.ai/p/post-one?_bhlid=abc123";
` + "```" + `
And this: https://codenewsletter.ai/p/post-two`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromCodeNewsletterURLs(input, &output)
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
