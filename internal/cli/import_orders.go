package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var importCmd = &cobra.Command{
	Use:   "import-orders",
	Short: "Import orders from JSON file",
}

func setupImportCmd(orderSvc order.Service) {
	importCmd.Flags().StringP(FlagFile, ShortFile, "", "File Path")

	_ = importCmd.MarkFlagRequired(FlagFile)

	importCmd.Run = func(cmd *cobra.Command, args []string) {
		path, err := getFlagString(cmd, FlagFile)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		count, err := orderSvc.ImportOrders(path)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		fmt.Printf("IMPORTED: %d\n", count)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
}
