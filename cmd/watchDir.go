package cmd

import (
	"github.com/gkwa/littlewill/watcher"
	"github.com/spf13/cobra"
)

var (
	patterns   []string
	filterType string
)

var watchDirCmd = &cobra.Command{
	Use:     "watch-dir [directory]",
	Aliases: []string{"wd"},
	Args:    cobra.ExactArgs(1),
	Short:   "Watch a directory for file changes",
	Long: `Watch a directory for file changes and process modified files.

You can specify patterns to filter which files to watch. If no patterns are specified,
all files will be watched. Patterns use standard glob syntax.

Examples:
  littlewill watch-dir /path/to/directory
  littlewill watch-dir /path/to/directory --patterns "*.md,*.txt"
  littlewill watch-dir /path/to/directory --patterns "doc_*.md" --patterns "report_*.txt"`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		watcher.RunWatcher(
			cmd.Context(),
			dir,
			patterns,
			filterType,
			linkTransforms,
		)
	},
}

func init() {
	rootCmd.AddCommand(watchDirCmd)

	watchDirCmd.Flags().StringSliceVarP(&patterns, "patterns", "p", []string{}, `File patterns to watch (comma-separated or multiple flags).
Examples:
  --patterns "*.md,*.txt"
  --patterns "*.go" --patterns "*.yaml"
Default: watch all files`)

	watchDirCmd.Flags().StringVarP(&filterType, "filter-type", "f", "write", `Event type to filter on. Options:
  create: Only watch for new files
  write: Watch for file modifications (default)
  remove: Watch for file deletions
  rename: Watch for file renames
  chmod: Watch for permission changes`)
}
