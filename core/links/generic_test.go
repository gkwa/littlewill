package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveGenericTrackingParams(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with UTM parameters",
			input:    "https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email&utm_campaign=2025-04-27-the-attack-you-invited",
			expected: "https://research.swtch.com/diffcover",
		},
		{
			name:     "URL with mixed parameters",
			input:    "https://example.com/article?id=123&utm_source=newsletter&page=2",
			expected: "https://example.com/article?id=123&page=2",
		},
		{
			name:     "URL without tracking parameters",
			input:    "https://example.com/page?id=123&category=tech",
			expected: "https://example.com/page?id=123&category=tech",
		},
		{
			name:     "URL with only tracking parameters",
			input:    "https://example.com/landing?utm_source=ads&utm_medium=social&utm_campaign=summer",
			expected: "https://example.com/landing",
		},
		{
			name:     "URL with Facebook click ID",
			input:    "https://example.com/product?id=456&fbclid=abc123",
			expected: "https://example.com/product?id=456",
		},
		{
			name: "Multiple URLs in text",
			input: `Check out these links:
https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email&utm_campaign=2025-04-27-the-attack-you-invited
https://example.com/article?id=123&utm_source=newsletter&page=2`,
			expected: `Check out these links:
https://research.swtch.com/diffcover
https://example.com/article?id=123&page=2`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email

` + "```" + `
var url = "https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email";
` + "```" + `

Another URL: https://example.com/article?id=123&utm_source=newsletter`,
			expected: `Check this URL: https://research.swtch.com/diffcover

` + "```" + `
var url = "https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email";
` + "```" + `

Another URL: https://example.com/article?id=123`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveGenericTrackingParams(input, &output)
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
