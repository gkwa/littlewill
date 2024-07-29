package watcher

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gkwa/littlewill/core"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var ignoredEvents = []fsnotify.Op{
	fsnotify.Chmod,
	fsnotify.Remove,
}

type EventHandler func(event fsnotify.Event, path string)

func RunWatcher(cmd *cobra.Command, args []string, patterns []string, filterType string, linkTransforms []func(io.Reader, io.Writer) error) {
	logger := logr.FromContextOrDiscard(cmd.Context())

	if len(args) == 0 {
		err := cmd.Usage()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
		}
		cmd.PrintErrln("Error: directory path is required")
		os.Exit(1)
	}

	dirToWatch := args[0]
	ctx := cmd.Context()
	handler := func(event fsnotify.Event, path string) {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Event: %s, File: %s\n", event.Op, path)
		err := core.ProcessFile(logger, path, linkTransforms...)
		if err != nil {
			logger.Error(err, "Failed to process file", "path", path)
		}
	}
	err := Run(ctx, dirToWatch, patterns, filterType, handler)
	if err != nil {
		cmd.PrintErrf("Error: %v\n", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, dirPath string, patterns []string, filterTypeStr string, handler EventHandler) error {
	logger := logr.FromContextOrDiscard(ctx)
	logger.Info("Starting directory watcher", "directory", dirPath)

	filters := makeFilters(patterns, parseFilterType(filterTypeStr))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if shouldTrigger(event, filters) && !isIgnoredEvent(event.Op) {
					absPath, err := filepath.Abs(event.Name)
					if err != nil {
						logger.Error(err, "Error getting absolute path", "file", event.Name)
						continue
					}
					handler(event, absPath)
					logger.V(1).Info("File event", "event", event.Op.String(), "file", absPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error(err, "Error from watcher")
			case <-ctx.Done():
				return
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		return fmt.Errorf("error adding directory to watcher: %w", err)
	}

	<-done
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
