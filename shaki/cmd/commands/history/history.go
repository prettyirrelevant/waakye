package history

import "github.com/spf13/cobra"

var HistoryCommand = cobra.Command{
	Use:   "history",
	Short: "Retrieve previous playlist conversions",
	Run:   func(cmd *cobra.Command, args []string) {},
}
