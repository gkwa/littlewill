package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveYouTubeCountFromLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "YouTube link with count",
			input:    "[(1619) Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft \"(1619) Understanding Neovim #1 - Installation - YouTube\")",
			expected: "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft \"(1619) Understanding Neovim #1 - Installation - YouTube\")",
		},
		{
			name:     "YouTube short link with count",
			input:    "[(2345) Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
			expected: "[Short YouTube Video](https://youtu.be/dQw4w9WgXcQ)",
		},
		{
			name:     "YouTube link without count",
			input:    "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft)",
			expected: "[Understanding Neovim #1 - Installation - YouTube](https://www.youtube.com/watch?v=87AXw9Quy9U&list=PLx2ksyallYzW4WNYHD9xOFrPRYGlntAft)",
		},
		{
			name:     "Non-YouTube link",
			input:    "[Example Link](https://example.com)",
			expected: "[Example Link](https://example.com)",
		},
		{
			name: "Multiple YouTube links",
			input: `
[(1234) First Video - YouTube](https://www.youtube.com/watch?v=abc123)
[(5678) Second Video - YouTube](https://youtu.be/def456 "Some title")
[Regular Link](https://example.com)
`,
			expected: `
[First Video - YouTube](https://www.youtube.com/watch?v=abc123)
[Second Video - YouTube](https://youtu.be/def456 "Some title")
[Regular Link](https://example.com)
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveYouTubeCountFromMarkdownLinks(input, &output)
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
