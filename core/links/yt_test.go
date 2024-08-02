package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromYouTubeLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "YouTube link with si and app parameters",
			input:    "https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq&app=Desktop",
			expected: "https://youtu.be/JSKJbGi5oNA",
		},
		{
			name:     "YouTube link without si parameter",
			input:    "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
		{
			name:     "Non-YouTube link",
			input:    "https://example.com?param=value",
			expected: "https://example.com?param=value",
		},
		{
			name: "Multiple YouTube links",
			input: `
Check out this video: https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq
And this one: https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=AnotherParam
`,
			expected: `
Check out this video: https://youtu.be/JSKJbGi5oNA
And this one: https://www.youtube.com/watch?v=dQw4w9WgXcQ
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveYoutubeParams(input, &output)
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
