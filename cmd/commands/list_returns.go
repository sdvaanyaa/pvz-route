package commands

import "github.com/spf13/cobra"

var listReturnsCmd = &cobra.Command{
	Use:   "list-returns",
	Short: "Get a list of returns",
}

func init() {
	rootCmd.AddCommand(listReturnsCmd)
}
