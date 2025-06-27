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
			name:     "URL with various UTM parameters",
			input:    "https://example.com/article?utm_source=newsletter&utm_medium=email&utm_campaign=summer&utm_content=button&utm_term=test&utm_id=123",
			expected: "https://example.com/article",
		},
		{
			name:     "URL with custom UTM parameters",
			input:    "https://example.com/page?utm_custom_param=value&utm_another_one=test&regular=keep",
			expected: "https://example.com/page?regular=keep",
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
			name:     "URL with Google Analytics parameters",
			input:    "https://example.com/page?_ga=abc123&_gl=def456&id=789",
			expected: "https://example.com/page?id=789",
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

func TestIsUTMParam(t *testing.T) {
	testCases := []struct {
		param    string
		expected bool
	}{
		{"utm_source", true},
		{"utm_medium", true},
		{"utm_campaign", true},
		{"utm_content", true},
		{"utm_term", true},
		{"utm_id", true},
		{"utm_custom_param", true},
		{"utm_whatever", true},
		{"utm_", true},
		{"source", false},
		{"medium", false},
		{"campaign", false},
		{"fbclid", false},
		{"gclid", false},
		{"utm", false},
		{"_utm_source", false},
	}

	for _, tc := range testCases {
		t.Run(tc.param, func(t *testing.T) {
			result := isUTMParam(tc.param)
			if result != tc.expected {
				t.Errorf("isUTMParam(%q) = %v, want %v", tc.param, result, tc.expected)
			}
		})
	}
}

func TestIsTrackingParam(t *testing.T) {
	testCases := []struct {
		param    string
		expected bool
	}{
		// UTM parameters
		{"utm_source", true},
		{"utm_medium", true},
		{"utm_campaign", true},
		{"utm_custom", true},
		// Common tracking parameters
		{"fbclid", true},
		{"gclid", true},
		{"_ga", true},
		{"source", true},
		{"medium", true},
		{"campaign", true},
		// Regular parameters that should be kept
		{"id", false},
		{"page", false},
		{"category", false},
		{"q", false},
		{"search", false},
	}

	for _, tc := range testCases {
		t.Run(tc.param, func(t *testing.T) {
			result := isTrackingParam(tc.param)
			if result != tc.expected {
				t.Errorf("isTrackingParam(%q) = %v, want %v", tc.param, result, tc.expected)
			}
		})
	}
}
