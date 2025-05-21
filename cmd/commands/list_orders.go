package commands

import "github.com/spf13/cobra"

var listOrdersCmd = &cobra.Command{
	Use:   "list-orders",
	Short: "Get a list of orders",
}

func init() {
	rootCmd.AddCommand(listOrdersCmd)
}
