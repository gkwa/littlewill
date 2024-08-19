package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestProcessPathsFromStdin(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Urls in code blocks with variable substitution are not adjusted",
			input: `` + "``` bash" + `
https://github.com/gkwa/${version}/test.bin` + "```" + ``,
			expected: `` + "``` bash" + `
https://github.com/gkwa/${version}/test.bin` + "```" + ``,
		},
		{
			name: "URLs outside code blocks are processed, URLs inside are not",
			input: `Check this link: https://www.google.com/search?q=test&hl=en
` + "```" + `
var url = "https://www.google.com/search?q=test&hl=en";
` + "```" + `
Another link: https://www.google.com/search?q=example&hl=en`,
			expected: `Check this link: https://www.google.com/search?q=test
` + "```" + `
var url = "https://www.google.com/search?q=test&hl=en";
` + "```" + `
Another link: https://www.google.com/search?q=example`,
		},
		{
			name: "Multiple code blocks with URLs are not processed",
			input: `Before code block
` + "```python" + `
url1 = "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=example"
` + "```" + `
Between code blocks
` + "```javascript" + `
const url2 = 'https://example.substack.com/p/article?utm_source=test';
` + "```" + `
After code block`,
			expected: `Before code block
` + "```python" + `
url1 = "https://www.youtube.com/watch?v=dQw4w9WgXcQ&si=example"
` + "```" + `
Between code blocks
` + "```javascript" + `
const url2 = 'https://example.substack.com/p/article?utm_source=test';
` + "```" + `
After code block`,
		},
		{
			name: "Nested code blocks are handled correctly",
			input: `Outer content
` + "```markdown" + `
# Markdown with code block

Here's a code block:

` + "```python" + `
url = "https://www.google.com/search?q=nested&hl=en"
` + "```" + `

End of markdown
` + "```" + `
Outer URL: https://www.google.com/search?q=outer&hl=en`,
			expected: `Outer content
` + "```markdown" + `
# Markdown with code block

Here's a code block:

` + "```python" + `
url = "https://www.google.com/search?q=nested&hl=en"
` + "```" + `

End of markdown
` + "```" + `
Outer URL: https://www.google.com/search?q=outer`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file to store the input content
			tempFile, err := os.CreateTemp("", "test_stdin_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			// Write the input content to the temporary file
			if _, err := tempFile.Write([]byte(tc.input)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tempFile.Close()

			// Create another temporary file to simulate stdin
			stdinFile, err := os.CreateTemp("", "test_stdin_paths_*.txt")
			if err != nil {
				t.Fatalf("Failed to create stdin temp file: %v", err)
			}
			defer os.Remove(stdinFile.Name())

			// Write the path of the input file to the stdin file
			if _, err := stdinFile.WriteString(tempFile.Name() + "\n"); err != nil {
				t.Fatalf("Failed to write to stdin temp file: %v", err)
			}
			stdinFile.Close()

			// Save the original stdin and restore it after the test
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Redirect stdin to read from the stdin file
			stdin, err := os.Open(stdinFile.Name())
			if err != nil {
				t.Fatalf("Failed to open stdin temp file: %v", err)
			}
			defer stdin.Close()
			os.Stdin = stdin

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create a test context
			ctx := context.Background()

			// Create a dummy cobra.Command
			cmd := &cobra.Command{}
			cmd.SetContext(ctx)

			// Run the pathsFromStdinCmd's Run function
			pathsFromStdinCmd.Run(cmd, []string{})

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("Failed to copy pipe to buffer: %v", err)
			}

			// Read the content of the processed file
			processedContent, err := os.ReadFile(tempFile.Name())
			if err != nil {
				t.Fatalf("Failed to read processed file: %v", err)
			}

			// Check the output
			if strings.TrimSpace(string(processedContent)) != tc.expected {
				t.Errorf("Expected output %q, but got %q", tc.expected, strings.TrimSpace(string(processedContent)))
			}
		})
	}
}
