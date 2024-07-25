//go:build windows
// +build windows

package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-logr/logr"
)

func lockFile(f *os.File) error {
	return nil
}

func unlockFile(f *os.File) error {
	return nil
}

func ProcessFile(ctx context.Context, path string, transform func(io.Reader, io.Writer) error) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger = logger.WithValues("file", path)
	logger.V(1).Info("Processing file")

	originalContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	var processedContent bytes.Buffer
	err = transform(bytes.NewReader(originalContent), &processedContent)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	if bytes.Equal(originalContent, processedContent.Bytes()) {
		logger.V(1).Info("File content unchanged, skipping write")
		return nil
	}

	originalFile, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open original file: %w", err)
	}
	defer func() {
		if err := unlockFile(originalFile); err != nil {
			logger.Error(err, "Failed to unlock file")
		}
		originalFile.Close()
	}()

	err = lockFile(originalFile)
	if err != nil {
		return err
	}

	_, err = originalFile.Write(processedContent.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write processed content to file: %w", err)
	}

	err = originalFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	logger.V(1).Info("Successfully processed and updated file")
	return nil
}

func ProcessPathsFromStdin(ctx context.Context, transform func(io.Reader, io.Writer) error) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.V(1).Info("Processing paths from stdin")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		path := scanner.Text()
		logger.V(1).Info("Processing path", "path", path)

		err := ProcessFile(ctx, path, transform)
		if err != nil {
			logger.Error(err, "Failed to process file", "path", path)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err, "Error reading input")
	}
}
