package links

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemoveParamsFromYouTubeURLs(t *testing.T) {
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
			err := RemoveParamsFromYouTubeURLs(input, &output)
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

func TestRemoveParamsFromSubstackURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Substack link with parameters",
			input:    "https://example.substack.com/p/article-title?utm_source=twitter&utm_medium=social",
			expected: "https://example.substack.com/p/article-title",
		},
		{
			name:     "Substack link without parameters",
			input:    "https://another.substack.com/p/another-article",
			expected: "https://another.substack.com/p/another-article",
		},
		{
			name:     "Non-Substack link",
			input:    "https://example.com?param=value",
			expected: "https://example.com?param=value",
		},
		{
			name: "Multiple Substack links",
			input: `
Check out this article: https://first.substack.com/p/title?utm_source=twitter
And this one: https://second.substack.com/p/another-title?utm_campaign=post
`,
			expected: `
Check out this article: https://first.substack.com/p/title
And this one: https://second.substack.com/p/another-title
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveParamsFromSubstackURLs(input, &output)
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

func TestRemoveTextFragments(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Link with text fragment",
			input:    "https://example.com/article#:~:text=some%20text",
			expected: "https://example.com/article",
		},
		{
			name:     "Link with regular fragment",
			input:    "https://example.com/article#heading-1",
			expected: "https://example.com/article#heading-1",
		},
		{
			name:     "Link without fragment",
			input:    "https://example.com/article",
			expected: "https://example.com/article",
		},
		{
			name: "Multiple links with text fragments",
			input: `
Check out this article: https://example.com/article1#:~:text=some%20text
And this one: https://example.com/article2#another-fragment
`,
			expected: `
Check out this article: https://example.com/article1
And this one: https://example.com/article2#another-fragment
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveTextFragments(input, &output)
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
