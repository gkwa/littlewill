package core

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
		{
			name: "Whitespace at beginning of line",
			input: `
			[

			Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
			
			https://substack.com/@javinpaul/note/c-62756371?r=21036
			
			https://chatgpt.com/share/650bd7e3-36ec-4e64-9503-953bfb09cf8b
`,
			expected: `
			[Managing 100s of Kubernetes Clusters using Cluster API](https://techblog.citystoragesystems.com/p/managing-100s-of-kubernetes-clusters "Managing 100s of Kubernetes Clusters using Cluster API")
			
			https://substack.com/@javinpaul/note/c-62756371?r=21036
			
			https://chatgpt.com/share/650bd7e3-36ec-4e64-9503-953bfb09cf8b
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
		{
			name: "Links with and without titles",
			input: `
[Google](https://google.com "Search Engine")
[GitHub](https://github.com)
[Stack Overflow](https://stackoverflow.com "Developer Community")
`,
			expected: `
[Google](https://google.com)
[GitHub](https://github.com)
[Stack Overflow](https://stackoverflow.com)
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

func TestRemoveTitlesFromMarkdownLinksErrors(t *testing.T) {
	t.Run("Read error", func(t *testing.T) {
		errReader := &errorReader{err: errors.New("read error")}
		var output bytes.Buffer
		err := RemoveTitlesFromMarkdownLinks(errReader, &output)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "RemoveTitlesFromMarkdownLinks: failed to read input: read error") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("Write error", func(t *testing.T) {
		input := strings.NewReader("some input")
		errWriter := &errorWriter{err: errors.New("write error")}
		err := RemoveTitlesFromMarkdownLinks(input, errWriter)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "RemoveTitlesFromMarkdownLinks: failed to write output: write error") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
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
