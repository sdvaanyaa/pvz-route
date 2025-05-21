package commands

import "github.com/spf13/cobra"

var issueCmd = &cobra.Command{
	Use:   "process-orders",
	Short: "Issue orders or accept customer returns",
}

func init() {
	rootCmd.AddCommand(issueCmd)
}
