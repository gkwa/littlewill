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

func ProcessFile(logger logr.Logger, path string, transforms ...func(io.Reader, io.Writer) error) error {
	logger.V(1).Info("Processing file", "path", path)

	originalContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	processedContent, err := ApplyTransforms(originalContent, transforms...)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	if bytes.Equal(originalContent, processedContent) {
		logger.V(1).Info("File content unchanged, skipping write")
		return nil
	}

	err = os.WriteFile(path, processedContent, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write processed content to file: %w", err)
	}

	logger.V(1).Info("Successfully processed and updated file")
	return nil
}

func ReadPathsFromStdin() ([]string, error) {
	var paths []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	return paths, scanner.Err()
}

func ProcessPaths(ctx context.Context, paths []string, transforms ...func(io.Reader, io.Writer) error) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.V(1).Info("Processing paths")

	for _, path := range paths {
		logger.V(1).Info("Processing path", "path", path)
		err := ProcessFile(logger, path, transforms...)
		if err != nil {
			logger.Error(err, "Failed to process file", "path", path)
		}
	}
}

func ProcessPathsFromStdin(ctx context.Context, transforms ...func(io.Reader, io.Writer) error) {
	logger := logr.FromContextOrDiscard(ctx)
	logger.V(1).Info("Processing paths from stdin")

	paths, err := ReadPathsFromStdin()
	if err != nil {
		logger.Error(err, "Error reading input")
		return
	}

	ProcessPaths(ctx, paths, transforms...)
}
