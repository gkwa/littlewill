package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileType(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "filetest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name     string
		setup    func() (string, error)
		expected string
	}{
		{
			name: "Text File",
			setup: func() (string, error) {
				path := filepath.Join(tmpDir, "test.txt")
				err := os.WriteFile(path, []byte("test"), 0o644)
				return path, err
			},
			expected: "Text File",
		},
		{
			name: "Go Source File",
			setup: func() (string, error) {
				path := filepath.Join(tmpDir, "test.go")
				err := os.WriteFile(path, []byte("package main"), 0o644)
				return path, err
			},
			expected: "Go Source File",
		},
		{
			name: "JPEG Image",
			setup: func() (string, error) {
				path := filepath.Join(tmpDir, "test.jpg")
				err := os.WriteFile(path, []byte("fake jpg"), 0o644)
				return path, err
			},
			expected: "JPEG Image",
		},
		{
			name: "Directory",
			setup: func() (string, error) {
				path := filepath.Join(tmpDir, "testdir")
				err := os.Mkdir(path, 0o755)
				return path, err
			},
			expected: "Directory",
		},
		{
			name: "Symlink",
			setup: func() (string, error) {
				target := filepath.Join(tmpDir, "target.txt")
				err := os.WriteFile(target, []byte("target"), 0o644)
				if err != nil {
					return "", err
				}
				link := filepath.Join(tmpDir, "symlink")
				err = os.Symlink(target, link)
				return link, err
			},
			expected: "Symlink to ", // We'll check if it starts with this
		},
		{
			name: "Unknown File Type",
			setup: func() (string, error) {
				path := filepath.Join(tmpDir, "test.xyz")
				err := os.WriteFile(path, []byte("unknown"), 0o644)
				return path, err
			},
			expected: "Unknown File Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			file := File{Path: path}
			result := file.FileType()

			if tt.name == "Symlink" {
				if !strings.HasPrefix(result, tt.expected) {
					t.Errorf("FileType() = %v, want prefix %v", result, tt.expected)
				}
			} else if result != tt.expected {
				t.Errorf("FileType() = %v, want %v", result, tt.expected)
			}
		})
	}
}
