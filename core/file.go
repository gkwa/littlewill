package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
)

func ProcessFile(ctx context.Context, path string) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger = logger.WithValues("file", path)
	logger.Info("Processing file")

	// Open the original file for reading
	originalFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open original file: %w", err)
	}
	defer originalFile.Close()

	// Create a temporary file in the same directory
	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, "littlewill-temp-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		// In case of any error, try to remove the temporary file
		os.Remove(tempPath)
	}()

	// Process the content
	err = CleanupMarkdownLinks(originalFile, tempFile)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	// Ensure all data is written to the temporary file
	err = tempFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync temporary file: %w", err)
	}

	// Close both files before renaming
	originalFile.Close()
	tempFile.Close()

	// Rename the temporary file to the original file
	err = os.Rename(tempPath, path)
	if err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	logger.Info("Successfully processed file")
	return nil
}
