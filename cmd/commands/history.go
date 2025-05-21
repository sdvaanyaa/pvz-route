package commands

import "github.com/spf13/cobra"

var historyCmd = &cobra.Command{
	Use:   "order-history",
	Short: "Get order change history",
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
