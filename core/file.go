package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
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

func ProcessFile(ctx context.Context, path string) error {
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

	dir := filepath.Dir(path)
	tempFile, err := os.CreateTemp(dir, "littlewill-temp-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(tempPath)
	}()

	_, err = originalFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to the beginning of original file: %w", err)
	}
	err = CleanupMarkdownLinks(originalFile, tempFile)
	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	err = tempFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync temporary file: %w", err)
	}

	tempFile.Close()

	err = os.Rename(tempPath, path)
	if err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	logger.Info("Successfully processed file")
	return nil
}

func RunPathsFromStdin(cmd *cobra.Command) {
	ctx := cmd.Context()
	logger := logr.FromContextOrDiscard(ctx)
	logger.V(1).Info("Processing paths from stdin")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		path := scanner.Text()
		logger.V(1).Info("Processing path", "path", path)

		err := ProcessFile(ctx, path)
		if err != nil {
			logger.Error(err, "Failed to process file", "path", path)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err, "Error reading input")
	}
}
