package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromBloombergURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Bloomberg article with accessToken and leadSource",
			input:    "https://www.bloomberg.com/news/articles/2025-12-03/anthropic-ceo-says-some-tech-firms-too-risky-with-ai-spending?accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb3VyY2UiOiJTdWJzY3JpYmVyR2lmdGVkQXJ0aWNsZSIsImlhdCI6MTc2NDgyOTI0OCwiZXhwIjoxNzY1NDM0MDQ4LCJhcnRpY2xlSWQiOiJUNlBJMkRUOTZPVDAwMCIsImJjb25uZWN0SWQiOiI2NTc1NjkyN0UwMkM0N0MwQkQ0MDNEQTJGMEUyNzIyMyJ9._HCGyu0Jjah9bwc0mehPz9S18cujW1D5hg2FWnZFVgo&leadSource=uverify+wall",
			expected: "https://www.bloomberg.com/news/articles/2025-12-03/anthropic-ceo-says-some-tech-firms-too-risky-with-ai-spending",
		},
		{
			name:     "Bloomberg article with only accessToken",
			input:    "https://www.bloomberg.com/news/articles/2025-12-03/some-article?accessToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9&other=keep",
			expected: "https://www.bloomberg.com/news/articles/2025-12-03/some-article?other=keep",
		},
		{
			name:     "Bloomberg article with only leadSource",
			input:    "https://www.bloomberg.com/news/articles/2025-12-03/some-article?leadSource=uverify+wall&keep=this",
			expected: "https://www.bloomberg.com/news/articles/2025-12-03/some-article?keep=this",
		},
		{
			name:     "Bloomberg article without tracking parameters",
			input:    "https://www.bloomberg.com/news/articles/2025-12-03/some-article?param=value",
			expected: "https://www.bloomberg.com/news/articles/2025-12-03/some-article?param=value",
		},
		{
			name:     "Non-Bloomberg URL with similar parameters",
			input:    "https://example.com/article?accessToken=test&leadSource=test",
			expected: "https://example.com/article?accessToken=test&leadSource=test",
		},
		{
			name: "Multiple Bloomberg URLs in text",
			input: `Check out these Bloomberg articles:
https://www.bloomberg.com/news/articles/2025-12-03/article1?accessToken=abc123&leadSource=email
https://www.bloomberg.com/news/articles/2025-12-04/article2?leadSource=twitter&other=keep`,
			expected: `Check out these Bloomberg articles:
https://www.bloomberg.com/news/articles/2025-12-03/article1
https://www.bloomberg.com/news/articles/2025-12-04/article2?other=keep`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Bloomberg URL: https://www.bloomberg.com/news/articles/2025-12-03/article?accessToken=test&leadSource=wall
` + "```" + `
var bloombergUrl = "https://www.bloomberg.com/news/articles/2025-12-03/article?accessToken=test&leadSource=wall";
` + "```" + `
Another Bloomberg URL: https://www.bloomberg.com/news/articles/2025-12-04/another?leadSource=email&accessToken=xyz`,
			expected: `Check this Bloomberg URL: https://www.bloomberg.com/news/articles/2025-12-03/article
` + "```" + `
var bloombergUrl = "https://www.bloomberg.com/news/articles/2025-12-03/article?accessToken=test&leadSource=wall";
` + "```" + `
Another Bloomberg URL: https://www.bloomberg.com/news/articles/2025-12-04/another`,
		},
		{
			name:     "Bloomberg subdomain URL",
			input:    "https://markets.bloomberg.com/data?accessToken=test123&leadSource=app&id=456",
			expected: "https://markets.bloomberg.com/data?id=456",
		},
		{
			name:     "Bloomberg URL with mixed parameters",
			input:    "https://www.bloomberg.com/news/articles/2025-12-03/article?id=123&accessToken=jwt_token&category=tech&leadSource=share&page=2",
			expected: "https://www.bloomberg.com/news/articles/2025-12-03/article?category=tech&id=123&page=2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromBloombergURLs(input, &output)
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
