package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromWSJURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "WSJ URL with reflink and st parameters",
			input:    "https://www.wsj.com/tech/waymo-cars-self-driving-robotaxi-tesla-uber-0777f570?reflink=desktopwebshare_permalink&st=p2QtzY",
			expected: "https://www.wsj.com/tech/waymo-cars-self-driving-robotaxi-tesla-uber-0777f570",
		},
		{
			name:     "WSJ URL with only reflink parameter",
			input:    "https://www.wsj.com/articles/some-article?reflink=desktopwebshare_permalink&other=keep",
			expected: "https://www.wsj.com/articles/some-article?other=keep",
		},
		{
			name:     "WSJ URL with only st parameter",
			input:    "https://www.wsj.com/business/article?st=abc123&param=value",
			expected: "https://www.wsj.com/business/article?param=value",
		},
		{
			name:     "WSJ URL without tracking parameters",
			input:    "https://www.wsj.com/articles/article-title?param=value",
			expected: "https://www.wsj.com/articles/article-title?param=value",
		},
		{
			name:     "Non-WSJ URL",
			input:    "https://example.com?reflink=something&st=tracking",
			expected: "https://example.com?reflink=something&st=tracking",
		},
		{
			name: "Multiple URLs in text",
			input: `Check out these links:
https://www.wsj.com/tech/article1?reflink=desktopwebshare_permalink&st=abc123
https://example.com/article?reflink=keep&st=also_keep`,
			expected: `Check out these links:
https://www.wsj.com/tech/article1
https://example.com/article?reflink=keep&st=also_keep`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://www.wsj.com/tech/article?reflink=desktopwebshare_permalink&st=abc123

` + "```" + `
var url = "https://www.wsj.com/tech/article?reflink=desktopwebshare_permalink&st=abc123";
` + "```" + `

Another URL: https://www.wsj.com/business/another?reflink=twitter&st=xyz789`,
			expected: `Check this URL: https://www.wsj.com/tech/article

` + "```" + `
var url = "https://www.wsj.com/tech/article?reflink=desktopwebshare_permalink&st=abc123";
` + "```" + `

Another URL: https://www.wsj.com/business/another`,
		},
		{
			name:     "WSJ URL with mixed parameters",
			input:    "https://www.wsj.com/markets/stocks/article?id=123&reflink=email&st=track123&section=business",
			expected: "https://www.wsj.com/markets/stocks/article?id=123&section=business",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromWSJURLs(input, &output)
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
