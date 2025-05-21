package commands

import "github.com/spf13/cobra"

var returnCmd = &cobra.Command{
	Use:   "return-order",
	Short: "Return the order to the courier",
}

func init() {
	rootCmd.AddCommand(returnCmd)
}
