package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromInstagramURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Instagram reel URL with igsh parameter",
			input:    "https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=NTc4MTIwNjQ2YQ%3D%3D",
			expected: "https://www.instagram.com/reel/DWnHqsCDlwC/",
		},
		{
			name:     "Instagram reel URL without tracking parameters",
			input:    "https://www.instagram.com/reel/DWnHqsCDlwC",
			expected: "https://www.instagram.com/reel/DWnHqsCDlwC",
		},
		{
			name:     "Instagram URL with igshid parameter",
			input:    "https://www.instagram.com/p/abc123/?igshid=xyz789",
			expected: "https://www.instagram.com/p/abc123/",
		},
		{
			name:     "Instagram URL with mixed parameters keeps non-tracking ones",
			input:    "https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=abc&other=keep",
			expected: "https://www.instagram.com/reel/DWnHqsCDlwC/?other=keep",
		},
		{
			name:     "Instagram profile URL with hl parameter",
			input:    "https://www.instagram.com/cannellevanille/?hl=en",
			expected: "https://www.instagram.com/cannellevanille/",
		},
		{
			name:     "Instagram search URL with %20-encoded space converts to plus",
			input:    "https://www.instagram.com/explore/search/keyword/?q=masienda%20tortilla%20press",
			expected: "https://www.instagram.com/explore/search/keyword/?q=masienda+tortilla+press",
		},
		{
			name:     "Instagram search URL with plus-encoded space unchanged",
			input:    "https://www.instagram.com/explore/search/keyword/?q=masiend+tortilla+press",
			expected: "https://www.instagram.com/explore/search/keyword/?q=masiend+tortilla+press",
		},
		{
			name:     "Non-Instagram URL is unchanged",
			input:    "https://example.com?igsh=something",
			expected: "https://example.com?igsh=something",
		},
		{
			name: "Multiple URLs in text",
			input: `Check these reels:
https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=NTc4MTIwNjQ2YQ%3D%3D
https://example.com/?igsh=keep`,
			expected: `Check these reels:
https://www.instagram.com/reel/DWnHqsCDlwC/
https://example.com/?igsh=keep`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=abc
` + "```" + `
var url = "https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=abc";
` + "```" + `
Another URL: https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=abc`,
			expected: `Check this URL: https://www.instagram.com/reel/DWnHqsCDlwC/
` + "```" + `
var url = "https://www.instagram.com/reel/DWnHqsCDlwC/?igsh=abc";
` + "```" + `
Another URL: https://www.instagram.com/reel/DWnHqsCDlwC/`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromInstagramURLs(input, &output)
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
