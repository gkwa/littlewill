package watcher

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type FilterType int

const (
	FilterTypeCreate FilterType = iota
	FilterTypeWrite
	FilterTypeRemove
	FilterTypeRename
	FilterTypeChmod
)

type Filter struct {
	Pattern string
	Type    FilterType
}

func shouldTrigger(event fsnotify.Event, filters []Filter) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		match, _ := filepath.Match(filter.Pattern, filepath.Base(event.Name))
		if match {
			switch filter.Type {
			case FilterTypeCreate:
				if event.Op&fsnotify.Create == fsnotify.Create {
					return true
				}
			case FilterTypeWrite:
				if event.Op&fsnotify.Write == fsnotify.Write {
					return true
				}
			case FilterTypeRemove:
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					return true
				}
			case FilterTypeRename:
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					return true
				}
			case FilterTypeChmod:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					return true
				}
			}
		}
	}
	return false
}
