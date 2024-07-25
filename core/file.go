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

func ProcessFile(ctx context.Context, path string, transforms ...func(io.Reader, io.Writer) error) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger = logger.WithValues("file", path)
	logger.V(1).Info("Processing file")

	originalContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	var processedContent bytes.Buffer
	currentContent := bytes.NewReader(originalContent)

	for _, transform := range transforms {
		processedContent.Reset()
		err = transform(currentContent, &processedContent)
		if err != nil {
			return fmt.Errorf("failed to process file: %w", err)
		}
		currentContent = bytes.NewReader(processedContent.Bytes())
	}

	if bytes.Equal(originalContent, processedContent.Bytes()) {
		logger.V(1).Info("File content unchanged, skipping write")
		return nil
	}

	err = os.WriteFile(path, processedContent.Bytes(), 0o644)
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
		err := ProcessFile(ctx, path, transforms...)
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
