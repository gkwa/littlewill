package cmd

import (
	"bufio"
	"os"

	"github.com/gkwa/littlewill/core"
	"github.com/spf13/cobra"
)

var pathsFromStdinCmd = &cobra.Command{
	Use:     "paths-from-stdin",
	Aliases: []string{"pfs"},
	Short:   "Process a list of paths from stdin",
	Long:    `This command reads a list of file paths from standard input and processes them, cleaning up markdown links in each file.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		logger := LoggerFrom(ctx)
		logger.V(1).Info("Processing paths from stdin")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			path := scanner.Text()
			logger.V(1).Info("Processing path", "path", path)

			err := core.ProcessFile(ctx, path)
			if err != nil {
				logger.Error(err, "Failed to process file", "path", path)
			}
		}

		if err := scanner.Err(); err != nil {
			logger.Error(err, "Error reading input")
		}
	},
}

func init() {
	rootCmd.AddCommand(pathsFromStdinCmd)
}
