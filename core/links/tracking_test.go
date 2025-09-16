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
			name:     "URL with fragment tracking parameters",
			input:    "https://musclemommasourdough.com/roasted-garlic-and-parmesan-sourdough-loaf/#growUnverifiedReaderId=4e75e57c-2848-4936-aa15-0e45cc7a0652&growAuthSource=espLink&growSource=espLink",
			expected: "https://musclemommasourdough.com/roasted-garlic-and-parmesan-sourdough-loaf/",
		},
		{
			name:     "URL with fragment tracking and regular fragment",
			input:    "https://example.com/page#section1&growAuthSource=espLink&utm_source=newsletter",
			expected: "https://example.com/page#section1",
		},
		{
			name:     "URL with mixed fragment parameters",
			input:    "https://example.com/page#keep=this&growAuthSource=espLink&other=also_keep&utm_source=newsletter",
			expected: "https://example.com/page#keep=this&other=also_keep",
		},
		{
			name:     "URL with regular fragment (no parameters)",
			input:    "https://example.com/page#heading-1",
			expected: "https://example.com/page#heading-1",
		},
		{
			name: "Multiple URLs in text",
			input: `Check out these links:
https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email&utm_campaign=2025-04-27-the-attack-you-invited
https://example.com/article?id=123&utm_source=newsletter&page=2
https://musclemommasourdough.com/recipe/#growUnverifiedReaderId=123&growAuthSource=espLink`,
			expected: `Check out these links:
https://research.swtch.com/diffcover
https://example.com/article?id=123&page=2
https://musclemommasourdough.com/recipe/`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email

` + "```" + `
var url = "https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email";
var fragmentUrl = "https://example.com/page#growAuthSource=espLink";
` + "```" + `

Another URL: https://example.com/article?id=123&utm_source=newsletter
Fragment URL: https://musclemommasourdough.com/recipe/#growUnverifiedReaderId=123`,
			expected: `Check this URL: https://research.swtch.com/diffcover

` + "```" + `
var url = "https://research.swtch.com/diffcover?utm_source=christophberger&utm_medium=email";
var fragmentUrl = "https://example.com/page#growAuthSource=espLink";
` + "```" + `

Another URL: https://example.com/article?id=123
Fragment URL: https://musclemommasourdough.com/recipe/`,
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
		// Fragment tracking parameters
		{"growUnverifiedReaderId", true},
		{"growAuthSource", true},
		{"growSource", true},
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

func TestParseFragmentParams(t *testing.T) {
	testCases := []struct {
		name     string
		fragment string
		expected map[string][]string
		hasError bool
	}{
		{
			name:     "Fragment with parameters",
			fragment: "growUnverifiedReaderId=123&growAuthSource=espLink&growSource=espLink",
			expected: map[string][]string{
				"growUnverifiedReaderId": {"123"},
				"growAuthSource":         {"espLink"},
				"growSource":             {"espLink"},
			},
			hasError: false,
		},
		{
			name:     "Fragment without parameters",
			fragment: "heading-1",
			expected: nil,
			hasError: false,
		},
		{
			name:     "Empty fragment",
			fragment: "",
			expected: nil,
			hasError: false,
		},
		{
			name:     "Fragment with mixed content",
			fragment: "section1&param=value&other=test",
			expected: map[string][]string{
				"section1": {""},
				"param":    {"value"},
				"other":    {"test"},
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseFragmentParams(tc.fragment)
			if tc.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tc.expected == nil && result != nil {
				t.Errorf("Expected nil result but got: %v", result)
			}
			if tc.expected != nil && result == nil {
				t.Errorf("Expected result but got nil")
			}
			if tc.expected != nil && result != nil {
				for key, expectedValues := range tc.expected {
					actualValues := result[key]
					if diff := cmp.Diff(expectedValues, actualValues); diff != "" {
						t.Errorf("Unexpected values for key %q (-want +got):\n%s", key, diff)
					}
				}
			}
		})
	}
}
