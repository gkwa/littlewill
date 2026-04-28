package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromTikTokURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "TikTok search URL with t timestamp parameter",
			input:    "https://www.tiktok.com/search?q=vegan+dumplings&t=1777392453994",
			expected: "https://www.tiktok.com/search?q=vegan+dumplings",
		},
		{
			name:     "TikTok search URL without tracking parameters is unchanged",
			input:    "https://www.tiktok.com/search?q=vegan+dumplings",
			expected: "https://www.tiktok.com/search?q=vegan+dumplings",
		},
		{
			name:     "TikTok URL with only t parameter",
			input:    "https://www.tiktok.com/@user/video/123?t=1777392453994",
			expected: "https://www.tiktok.com/@user/video/123",
		},
		{
			name:     "Non-TikTok URL with t parameter is unchanged",
			input:    "https://www.youtube.com/watch?v=abc&t=42",
			expected: "https://www.youtube.com/watch?v=abc&t=42",
		},
		{
			name: "Multiple URLs in text",
			input: `Check these:
https://www.tiktok.com/search?q=vegan+dumplings&t=1777392453994
https://www.tiktok.com/search?q=vegan+dumplings`,
			expected: `Check these:
https://www.tiktok.com/search?q=vegan+dumplings
https://www.tiktok.com/search?q=vegan+dumplings`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromTikTokURLs(input, &output)
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
