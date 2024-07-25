package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"

	"github.com/go-logr/logr"
)

var (
	lockFile   func(*os.File) error
	unlockFile func(*os.File) error
)

func init() {
	if runtime.GOOS == "windows" {
		lockFile = func(f *os.File) error {
			return nil
		}
		unlockFile = func(f *os.File) error {
			return nil
		}
	} else {
		lockFile = func(f *os.File) error {
			err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
			if err != nil {
				return fmt.Errorf("failed to lock file: %w", err)
			}
			return nil
		}
		unlockFile = func(f *os.File) error {
			err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
			if err != nil {
				return fmt.Errorf("failed to unlock file: %w", err)
			}
			return nil
		}
	}
}

func ProcessFile(ctx context.Context, path string, transform func(io.Reader, io.Writer) error) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger = logger.WithValues("file", path)
	logger.Info("Processing file")

	originalFile, err := os.OpenFile(path, os.O_RDWR, 0o644)
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

	originalContent, err := io.ReadAll(originalFile)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	var processedContent bytes.Buffer
	err = transform(bytes.NewReader(originalContent), &processedContent)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	if bytes.Equal(originalContent, processedContent.Bytes()) {
		logger.Info("File content unchanged, skipping write")
		return nil
	}

	_, err = originalFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to the beginning of original file: %w", err)
	}

	err = originalFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate original file: %w", err)
	}

	_, err = io.Copy(originalFile, &processedContent)
	if err != nil {
		return fmt.Errorf("failed to write processed content to file: %w", err)
	}

	err = originalFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	logger.Info("Successfully processed and updated file")
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
