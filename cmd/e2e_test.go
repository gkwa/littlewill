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
			name: "Google Search Link",
			input: `
			
			[  

qfc Chocolate Ice Cream - Google Search](https://www.google.com/search?q=qfc+Chocolate+Ice+Cream&oq=qfc+Chocolate+Ice+Cream&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIGCAEQRRhA0gEIMTU2NmowajSoAgCwAgE&sourceid=chrome&ie=UTF-8 "qfc Chocolate Ice Cream - Google Search")

`,
			expected: `[qfc Chocolate Ice Cream - Google Search](https://www.google.com/search?q=qfc+Chocolate+Ice+Cream)`,
		},

		{
			name: "Copying text into obsidian brings in funy formatting that breaks links",
			input: `



			![ 
			
			Item media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)
			

`,
			expected: `![Item media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)`,
		},
		{
			name: "Copying text into obsidian brings in funy formatting that breaks links - part 2",
			input: `



			[Item

			media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)
			

`,
			expected: `[Item media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)`,
		},
		{
			name: "Copying text into obsidian brings in funy formatting that breaks links - part 3",
			input: `



			[ Item

			media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)
			

`,
			expected: `[Item media 5 screenshot](https://lh3.googleusercontent.com/ADTbiH2FM2SYb3PbxWeAI0v_-FYVMFt_6hJ3sabl_gVDadugPc5FX55USRMRIo50uvD0gwKqIJu-kfXWJHRiQV6SsTE=s1280-w1280-h800)`,
		},
		{
			name: "Urls in code blocks with variable substitution are not adjusted",
			input: `` + "``` bash" + `
https://github.com/gkwa/${version}/test.bin` + "```" + ``,
			expected: `` + "``` bash" + `
https://github.com/gkwa/${version}/test.bin` + "```" + ``,
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
