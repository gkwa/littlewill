package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromGoogleLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Google search link with parameters",
			input:    "https://www.google.com/search?q=test&bih=722&biw=1536&hl=en&sxsrf=ABC123",
			expected: "https://www.google.com/search?hl=en&q=test",
		},
		{
			name:     "Google link without parameters",
			input:    "https://www.google.com",
			expected: "https://www.google.com",
		},
		{
			name:     "Google Maps link (excluded)",
			input:    "https://www.google.com/maps/place/New+York",
			expected: "https://www.google.com/maps/place/New+York",
		},
		{
			name:     "Non-Google link",
			input:    "https://example.com?param=value",
			expected: "https://example.com?param=value",
		},
		{
			name: "Multiple Google links",
			input: `
Search result: https://www.google.com/search?q=test&bih=722&biw=1536&hl=en&sxsrf=ABC123
Maps link: https://www.google.com/maps/place/New+York
`,
			expected: `
Search result: https://www.google.com/search?hl=en&q=test
Maps link: https://www.google.com/maps/place/New+York
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromGoogleURLs(input, &output)
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
