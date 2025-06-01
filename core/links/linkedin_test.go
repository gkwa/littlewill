package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromLinkedInURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "LinkedIn post with rcm parameter",
			input:    "https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/?rcm=ACoAABMsKDgBh79fm4Bw4lyI4TyVv1ADaplZsFA",
			expected: "https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/",
		},
		{
			name:     "LinkedIn post with multiple tracking parameters",
			input:    "https://www.linkedin.com/posts/johndoe_software-engineering-activity-123456789?rcm=ABC123&utm_source=newsletter&utm_medium=email",
			expected: "https://www.linkedin.com/posts/johndoe_software-engineering-activity-123456789",
		},
		{
			name:     "LinkedIn URL without tracking parameters",
			input:    "https://www.linkedin.com/in/johndoe",
			expected: "https://www.linkedin.com/in/johndoe",
		},
		{
			name:     "Non-LinkedIn URL",
			input:    "https://example.com?rcm=shouldnotberemoved",
			expected: "https://example.com?rcm=shouldnotberemoved",
		},
		{
			name: "Multiple LinkedIn URLs in text",
			input: `Check out these LinkedIn posts:
https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/?rcm=ACoAABMsKDgBh79fm4Bw4lyI4TyVv1ADaplZsFA
https://www.linkedin.com/posts/johndoe_software-engineering-activity-123456789?utm_source=newsletter`,
			expected: `Check out these LinkedIn posts:
https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/
https://www.linkedin.com/posts/johndoe_software-engineering-activity-123456789`,
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this LinkedIn URL: https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/?rcm=ACoAABMsKDgBh79fm4Bw4lyI4TyVv1ADaplZsFA

` + "```" + `
var linkedinUrl = "https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/?rcm=ACoAABMsKDgBh79fm4Bw4lyI4TyVv1ADaplZsFA";
` + "```" + `

Another LinkedIn URL: https://www.linkedin.com/posts/johndoe_activity-123456789?rcm=TEST123`,
			expected: `Check this LinkedIn URL: https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/

` + "```" + `
var linkedinUrl = "https://www.linkedin.com/posts/nikkisiapno_api-gateway-vs-load-balancer-whats-the-activity-7332371450495963137-tQDq/?rcm=ACoAABMsKDgBh79fm4Bw4lyI4TyVv1ADaplZsFA";
` + "```" + `

Another LinkedIn URL: https://www.linkedin.com/posts/johndoe_activity-123456789`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromLinkedInURLs(input, &output)
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
