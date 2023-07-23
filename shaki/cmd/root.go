package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/prettyirrelevant/shaki/cmd/commands/convert"
)

const Version = "0.0.4"

var rootCmd = &cobra.Command{
	Use:     "waakye",
	Short:   "This is the CLI application for the playlist converter, waakye.",
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(convert.ConvertCmd)
}
