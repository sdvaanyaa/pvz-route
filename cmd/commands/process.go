package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
	"strings"
)

var processCmd = &cobra.Command{
	Use:   "process-orders",
	Short: "Issue orders or accept customer returns",
}

func SetupProcessCmd(orderService *services.OrderService) {
	processCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	processCmd.Flags().StringP(FlagAction, "a", "", "Expires")
	processCmd.Flags().StringP(FlagOrderIDs, "o", "", "Order ID")

	_ = processCmd.MarkFlagRequired(FlagOrderID)
	_ = processCmd.MarkFlagRequired(FlagUserID)
	_ = processCmd.MarkFlagRequired(FlagExpires)

	processCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, err := GetFlag(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		action, err := GetFlag(cmd, FlagAction)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		orderIDs, err := GetFlag(cmd, FlagOrderIDs)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		sliceOrderIDs := strings.Split(orderIDs, ",")

		for _, orderID := range sliceOrderIDs {
			if err = orderService.ProcessOrder(userID, orderID, action); err != nil {
				fmt.Printf("ERROR %s: %s\n", orderID, err)
				continue
			}

			fmt.Printf("PROCESSED: %s\n", orderID)

		}
	}
}

func init() {
	rootCmd.AddCommand(processCmd)
}
