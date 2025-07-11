package watcher

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gkwa/littlewill/core"
	"github.com/go-logr/logr"
)

var ignoredEvents = []fsnotify.Op{
	fsnotify.Chmod,
	fsnotify.Remove,
	fsnotify.Rename,
}

type EventHandler func(event fsnotify.Event, path string)

func isBrokenSymlink(path string) bool {
	// Lstat gets info about the symlink itself (doesn't follow it)
	linkInfo, err := os.Lstat(path)
	if err != nil || linkInfo.Mode()&os.ModeSymlink == 0 {
		return false // Not a symlink or can't read it
	}

	// Stat follows the symlink to the target
	_, err = os.Stat(path)
	return err != nil // If Stat fails, the target doesn't exist
}

func cleanupBrokenSymlinks(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't read
		}

		if isBrokenSymlink(path) {
			fmt.Printf("Removing broken symlink: %s\n", path)
			return os.Remove(path)
		}
		return nil
	})
}

func addWatchWithCleanup(w *fsnotify.Watcher, dirPath string) error {
	err := w.Add(dirPath)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		// Extract the problematic file path from the error message
		if strings.Contains(err.Error(), ".#") {
			fmt.Printf("Broken symlink detected, cleaning up and retrying...\n")
			cleanupBrokenSymlinks(dirPath)

			// Retry after cleanup
			time.Sleep(1000 * time.Millisecond)
			return w.Add(dirPath)
		}
	}
	return err
}

func RunWatcher(
	ctx context.Context,
	dirToWatch string,
	patterns []string,
	filterType string,
	linkTransforms []func(io.Reader, io.Writer) error,
) {
	logger := logr.FromContextOrDiscard(ctx)

	handler := func(event fsnotify.Event, path string) {
		time.Sleep(time.Second)
		fmt.Printf("Event: %s, File: %s\n", event.Op, path)

		err := core.ProcessFile(logger, path, linkTransforms...)
		if err != nil {
			logger.Error(err, "Failed to process file", "path", path)
			// Don't exit on file processing errors, just continue watching
		}
	}
	err := Run(ctx, dirToWatch, patterns, filterType, handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func Run(
	ctx context.Context,
	dirPath string,
	patterns []string,
	filterTypeStr string,
	handler EventHandler,
) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("Starting directory watcher", "directory", dirPath)

	filters := makeFilters(patterns, parseFilterType(filterTypeStr))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logger.Info("Watcher events channel closed")
					return
				}
				if shouldTrigger(event, filters) && !isIgnoredEvent(event.Op) {
					absPath, err := filepath.Abs(event.Name)
					if err != nil {
						logger.Error(err, "Error getting absolute path", "file", event.Name)
						continue
					}
					// Handle the event in a separate goroutine to prevent blocking
					go func(evt fsnotify.Event, path string) {
						defer func() {
							if r := recover(); r != nil {
								logger.Error(fmt.Errorf("panic in event handler: %v", r), "Recovered from panic", "file", path)
							}
						}()
						handler(evt, path)
					}(event, absPath)
					logger.V(1).Info("File event", "event", event.Op.String(), "file", absPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logger.Info("Watcher errors channel closed")
					return
				}
				logger.Error(err, "Error from watcher")
				// Continue watching even if there are watcher errors
			case <-ctx.Done():
				logger.Info("Context cancelled, stopping watcher")
				return
			}
		}
	}()

	err = addWatchWithCleanup(watcher, dirPath)
	if err != nil {
		return fmt.Errorf("error adding directory to watcher: %w", err)
	}

	logger.Info("Watcher started successfully", "directory", dirPath)

	// Block until context is cancelled
	<-ctx.Done()
	logger.Info("Watcher stopping")
	return nil
}

func isIgnoredEvent(op fsnotify.Op) bool {
	for _, ignoredOp := range ignoredEvents {
		if op == ignoredOp {
			return true
		}
	}
	return false
}

func makeFilters(patterns []string, filterType FilterType) []Filter {
	filters := make([]Filter, 0, len(patterns))
	for _, pattern := range patterns {
		filters = append(filters, Filter{
			Pattern: pattern,
			Type:    filterType,
		})
	}
	return filters
}

func parseFilterType(filterType string) FilterType {
	switch filterType {
	case "create":
		return FilterTypeCreate
	case "write":
		return FilterTypeWrite
	case "remove":
		return FilterTypeRemove
	case "rename":
		return FilterTypeRename
	case "chmod":
		return FilterTypeChmod
	default:
		return FilterTypeWrite
	}
}
