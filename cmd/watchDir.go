package cmd

import (
	"os"

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
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Usage()
			if err != nil {
				cmd.PrintErrf("Error: %v\n", err)
			}
			cmd.PrintErrln("Error: directory path is required")
			os.Exit(1)
		}

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
	watchDirCmd.Flags().StringSliceVar(&patterns, "patterns", []string{}, "File patterns to watch")
	watchDirCmd.Flags().StringVar(&filterType, "filter-type", "write", "Filter type (create, write, remove, rename, chmod)")
}
