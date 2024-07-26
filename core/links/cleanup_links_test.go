package links

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanupMarkdownLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "No links",
			input: `
This is some text.
This more text.
Here is a link https://google.com.
`,
			expected: `
This is some text.
This more text.
Here is a link https://google.com.
`,
		},
		{
			name: "Single link with extra space",
			input: `
[    Google](https://google.com)
This is some text
[

Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
This more text.
Here is a link https://google.com.
`,
			expected: `
[Google](https://google.com)
This is some text
[Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
This more text.
Here is a link https://google.com.
`,
		},
		{
			name: "Multiple links with extra space",
			input: `
[Google](https://google.com)
This is some text
[
Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
This more text.
[

  Another link with space](https://example.com)
Here is a link https://google.com.
`,
			expected: `
[Google](https://google.com)
This is some text
[Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
This more text.
[Another link with space](https://example.com)
Here is a link https://google.com.
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveWhitespaceFromMarkdownLinks(input, &output)
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

func TestCleanupMarkdownLinksErrors(t *testing.T) {
	t.Run("Read error", func(t *testing.T) {
		errReader := &errorReader{err: errors.New("read error")}
		var output bytes.Buffer
		err := RemoveWhitespaceFromMarkdownLinks(errReader, &output)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "CleanupMarkdownLinks: failed to read input: read error") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Write error", func(t *testing.T) {
		input := strings.NewReader("some input")
		errWriter := &errorWriter{err: errors.New("write error")}
		err := RemoveWhitespaceFromMarkdownLinks(input, errWriter)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "CleanupMarkdownLinks: failed to write output: write error") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
}

func TestRemoveTitlesFromMarkdownLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "No links",
			input: `
This is some text.
This more text.
Here is a link https://google.com.
`,
			expected: `
This is some text.
This more text.
Here is a link https://google.com.
`,
		},
		{
			name:     "Link with title",
			input:    `[Google](https://google.com "Search Engine")`,
			expected: `[Google](https://google.com)`,
		},
		{
			name: "Multiple links with titles",
			input: `
[Google](https://google.com "Search Engine")
This is some text
[Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
This more text.
[Another link](https://example.com "Example")
Here is a link https://google.com.
`,
			expected: `
[Google](https://google.com)
This is some text
[Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters)
This more text.
[Another link](https://example.com)
Here is a link https://google.com.
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			var output bytes.Buffer
			err := RemoveTitlesFromMarkdownLinks(input, &output)
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

func TestRemoveParamsFromYouTubeLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "YouTube link with si parameter",
			input:    "https://youtu.be/JSKJbGi5oNA?si=b2GkFDivckm1k-Mq",
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

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
}
