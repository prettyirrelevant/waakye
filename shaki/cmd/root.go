package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/prettyirrelevant/shaki/cmd/commands/convert"
	"github.com/prettyirrelevant/shaki/cmd/commands/history"
)

var rootCmd = &cobra.Command{
	Use:   "shaki",
	Short: "shaki is the CLI application for Waakye",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convert.ConvertCmd)
	rootCmd.AddCommand(&history.HistoryCommand)
}
