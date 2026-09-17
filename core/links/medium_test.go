package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromMediumURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Medium byline URL with source parameter",
			input:    "https://medium.com/@albertogeniola?source=post_page---byline--444ffe641b64--------------------------------",
			expected: "https://medium.com/@albertogeniola",
		},
		{
			name:     "Medium subdomain URL with source parameter",
			input:    "https://blog.medium.com/some-post?source=rss&id=1",
			expected: "https://blog.medium.com/some-post?id=1",
		},
		{
			name:     "Medium URL without tracking parameters",
			input:    "https://medium.com/@someone/post?id=1",
			expected: "https://medium.com/@someone/post?id=1",
		},
		{
			name:     "Lookalike host keeps source parameter",
			input:    "https://notmedium.com/page?source=first",
			expected: "https://notmedium.com/page?source=first",
		},
		{
			name:     "OSRM trip URL keeps source parameter",
			input:    "https://routing.openstreetmap.de/routed-foot/trip/v1/foot/-122.3,47.6;-122.2,47.5?source=first&destination=any&roundtrip=false",
			expected: "https://routing.openstreetmap.de/routed-foot/trip/v1/foot/-122.3,47.6;-122.2,47.5?source=first&destination=any&roundtrip=false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromMediumURLs(input, &output)
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
