package cmd

import (
	"github.com/gkwa/littlewill/core"
	"github.com/spf13/cobra"
)

var pathsFromStdinCmd = &cobra.Command{
	Use:     "paths-from-stdin",
	Aliases: []string{"pfs"},
	Short:   "Process a list of paths from stdin",
	Long:    `This command reads a list of file paths from standard input and processes them, cleaning up markdown links in each file.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.ProcessPathsFromStdin(
			cmd.Context(),
			core.RemoveWhitespaceFromMarkdownLinks,
			core.RemoveTitlesFromMarkdownLinks,
			core.RemoveParamsFromYouTubeLinks,
		)
	},
}

func init() {
	rootCmd.AddCommand(pathsFromStdinCmd)
}
