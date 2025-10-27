package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromMailchimpURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Mailchimp redirect URL with e parameter",
			input:    "https://mailchi.mp/quantamagazine.org/why-black-holes-keep-pulling-physicists-in-4867774?e=8df2929e90",
			expected: "https://mailchi.mp/quantamagazine.org/why-black-holes-keep-pulling-physicists-in-4867774",
		},
		{
			name:     "Mailchimp URL with e parameter and other params",
			input:    "https://mailchi.mp/example.com/article?e=abc123def456&id=123&title=test",
			expected: "https://mailchi.mp/example.com/article?id=123&title=test",
		},
		{
			name:     "Mailchimp URL with utm and e parameters",
			input:    "https://mailchi.mp/newsletter.com/feature?e=xyz789&utm_source=mailchimp&utm_medium=email&utm_campaign=weekly",
			expected: "https://mailchi.mp/newsletter.com/feature",
		},
		{
			name:     "Mailchimp URL without tracking parameters",
			input:    "https://mailchi.mp/example.com/page?id=456&category=news",
			expected: "https://mailchi.mp/example.com/page?category=news&id=456",
		},
		{
			name:     "Mailchimp URL with only e parameter",
			input:    "https://mailchi.mp/blog.com/post?e=subscriber123",
			expected: "https://mailchi.mp/blog.com/post",
		},
		{
			name:     "Non-Mailchimp URL with e parameter",
			input:    "https://example.com?e=shouldnotberemoved",
			expected: "https://example.com?e=shouldnotberemoved",
		},
		{
			name:     "Mailchimp domain URL",
			input:    "https://mailchimp.com/campaign?e=track123&utm_source=internal",
			expected: "https://mailchimp.com/campaign",
		},
		{
			name: "Multiple Mailchimp URLs in text",
			input: `Check out these links:
https://mailchi.mp/quantamagazine.org/article?e=8df2929e90
https://mailchi.mp/techblog.com/post?e=abc123&id=456&utm_source=newsletter
https://example.com/page?e=keep&title=test`,
			expected: `Check out these links:
https://mailchi.mp/quantamagazine.org/article
https://mailchi.mp/techblog.com/post?id=456&title=test
https://example.com/page?e=keep&title=test`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this Mailchimp URL: https://mailchi.mp/quantamagazine.org/article?e=8df2929e90&utm_source=email

` + "```" + `
var mailchimpUrl = "https://mailchi.mp/example.com/page?e=subscriber123&utm_source=mailchimp";
` + "```" + `

Another Mailchimp URL: https://mailchi.mp/blog.com/post?e=track456&id=789`,
			expected: `Check this Mailchimp URL: https://mailchi.mp/quantamagazine.org/article

` + "```" + `
var mailchimpUrl = "https://mailchi.mp/example.com/page?e=subscriber123&utm_source=mailchimp";
` + "```" + `

Another Mailchimp URL: https://mailchi.mp/blog.com/post?id=789`,
		},
		{
			name:     "Mailchimp URL with multiple e parameters (edge case)",
			input:    "https://mailchi.mp/example.com/article?e=first&e=second&keep=value",
			expected: "https://mailchi.mp/example.com/article?keep=value",
		},
		{
			name:     "Mailchimp subdomain URL",
			input:    "https://mail.mailchi.mp/campaign/feature?e=sub123&utm_campaign=monthly",
			expected: "https://mail.mailchi.mp/campaign/feature",
		},
		{
			name:     "Mailchimp URL with utm_source and e parameter",
			input:    "https://mailchi.mp/newsletter.org/issue-42?e=reader789&utm_source=mailchimp&sort=recent",
			expected: "https://mailchi.mp/newsletter.org/issue-42?sort=recent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromMailchimpURLs(input, &output)
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

func TestIsMailchimpURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected bool
	}{
		{"https://mailchi.mp/example.com/page", true},
		{"https://mail.mailchi.mp/campaign", true},
		{"https://mailchimp.com/features", true},
		{"https://app.mailchimp.com/dashboard", true},
		{"https://example.com/page", false},
		{"https://notmailchimp.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			u, _ := url.Parse(tc.url)
			result := isMailchimpURL(u)
			if result != tc.expected {
				t.Errorf("isMailchimpURL(%q) = %v, want %v", tc.url, result, tc.expected)
			}
		})
	}
}

func TestIsMailchimpTrackingParam(t *testing.T) {
	testCases := []struct {
		param    string
		expected bool
	}{
		// Mailchimp-specific parameters
		{"e", true},
		// UTM parameters (shared logic)
		{"utm_source", true},
		{"utm_medium", true},
		{"utm_campaign", true},
		// Regular parameters that should be kept
		{"id", false},
		{"page", false},
		{"title", false},
		{"category", false},
	}

	for _, tc := range testCases {
		t.Run(tc.param, func(t *testing.T) {
			result := isMailchimpTrackingParam(tc.param)
			if result != tc.expected {
				t.Errorf("isMailchimpTrackingParam(%q) = %v, want %v", tc.param, result, tc.expected)
			}
		})
	}
}
