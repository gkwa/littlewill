package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gkwa/littlewill/core"
	"github.com/gkwa/littlewill/core/links"
	"github.com/gkwa/littlewill/watcher"
	"github.com/spf13/cobra"
)

var (
	dirToWatch string
	patterns   []string
	filterType string
)

// watchDirCmd represents the watchDir command
var watchDirCmd = &cobra.Command{
	Use:     "watch-dir",
	Aliases: []string{"wd"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())

		if len(args) == 0 {
			err := cmd.Usage()
			if err != nil {
				cmd.PrintErrf("Error: %v\n", err)
			}
			cmd.PrintErrln("Error: directory path is required")
			os.Exit(1)
		}

		dirToWatch = args[0]
		ctx := cmd.Context()
		handler := func(event fsnotify.Event, path string) {
			time.Sleep(500 * time.Millisecond)

			fmt.Printf("Event: %s, File: %s\n", event.Op, path)
			err := core.ProcessFile(
				logger,
				path,
				links.RemoveWhitespaceFromMarkdownLinks,
				links.RemoveTitlesFromMarkdownLinks,
				links.RemoveParamsFromYouTubeURLs,
				links.RemoveParamsFromGoogleURLs,
				links.RemoveYouTubeCountFromMarkdownLinks,
			)
			if err != nil {
				logger.Error(err, "Failed to process file", "path", path)
			}
		}
		err := watcher.Run(ctx, dirToWatch, patterns, filterType, handler)
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchDirCmd)
}
