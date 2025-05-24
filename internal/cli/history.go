package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var historyCmd = &cobra.Command{
	Use:   "order-history",
	Short: "Get order change history",
}

func setupHistoryCmd(orderSvc order.Service) {
	historyCmd.Run = func(cmd *cobra.Command, args []string) {
		historyEntry, err := orderSvc.History()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		for _, entry := range historyEntry {
			formatedTime := entry.Timestamp.Format("02 Jan 15:04")
			fmt.Printf("HISTORY: %s %s %s\n", entry.OrderID, entry.Status, formatedTime)
		}
	}
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
