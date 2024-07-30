package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path string
}

func (f File) IsSymlink() (bool, error) {
	fileInfo, err := os.Lstat(f.Path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

func (f File) FileType() string {
	// Get file info without following symlinks
	info, err := os.Lstat(f.Path)
	if err != nil {
		return "Error: Unable to get file info"
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		// If it's a symlink, we can optionally get the target
		target, err := os.Readlink(f.Path)
		if err != nil {
			return "Symlink (unable to read target)"
		}
		return fmt.Sprintf("Symlink to %s", target)
	}

	// If it's not a symlink, proceed with file type detection
	if info.IsDir() {
		return "Directory"
	}

	ext := strings.ToLower(filepath.Ext(f.Path))
	switch ext {
	case ".txt":
		return "Text File"
	case ".go":
		return "Go Source File"
	case ".jpg", ".jpeg":
		return "JPEG Image"
	case ".png":
		return "PNG Image"
	case ".pdf":
		return "PDF Document"
	default:
		return "Unknown File Type"
	}
}
