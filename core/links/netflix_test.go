package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromNetflixURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Netflix URL without tracking parameters",
			input:    "https://www.netflix.com/us/title/70206131",
			expected: "https://www.netflix.com/us/title/70206131",
		},
		{
			name:     "Netflix URL with tracking parameters",
			input:    "https://www.netflix.com/us/title/70206131?s=a&trkid=13747225&shareType=Title&shareUuid=b2108556-c03f-48bf-b18c-71aa714ce7c0&trg=more&unifiedEntityIdEncoded=Video%3A70206131&vlang=en&clip=81375738",
			expected: "https://www.netflix.com/us/title/70206131",
		},
		{
			name:     "Netflix URL with partial tracking parameters",
			input:    "https://www.netflix.com/us/title/70206131?trkid=13747225&vlang=en",
			expected: "https://www.netflix.com/us/title/70206131",
		},
		{
			name:     "Netflix URL with mixed parameters keeps non-tracking ones",
			input:    "https://www.netflix.com/us/title/70206131?trkid=13747225&other=keep",
			expected: "https://www.netflix.com/us/title/70206131?other=keep",
		},
		{
			name:     "Non-Netflix URL is unchanged",
			input:    "https://example.com?trkid=something&shareType=Title",
			expected: "https://example.com?trkid=something&shareType=Title",
		},
		{
			name: "Multiple URLs in text",
			input: `Watch these shows:
https://www.netflix.com/us/title/70206131?s=a&trkid=13747225&shareType=Title
https://example.com/title/70206131?trkid=keep`,
			expected: `Watch these shows:
https://www.netflix.com/us/title/70206131
https://example.com/title/70206131?trkid=keep`,
		},
		{
			name:     "Netflix URL with trailing slash gets it stripped",
			input:    "https://www.netflix.com/us/title/70206131/",
			expected: "https://www.netflix.com/us/title/70206131",
		},
		{
			name: "URLs inside code blocks should not be processed",
			input: `Check this URL: https://www.netflix.com/us/title/70206131?s=a&trkid=13747225
` + "```" + `
var url = "https://www.netflix.com/us/title/70206131?s=a&trkid=13747225";
` + "```" + `
Another URL: https://www.netflix.com/us/title/70206131?s=a&trkid=13747225`,
			expected: `Check this URL: https://www.netflix.com/us/title/70206131
` + "```" + `
var url = "https://www.netflix.com/us/title/70206131?s=a&trkid=13747225";
` + "```" + `
Another URL: https://www.netflix.com/us/title/70206131`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer

			err := RemoveParamsFromNetflixURLs(input, &output)
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
