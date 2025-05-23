package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
	"strings"
)

var processCmd = &cobra.Command{
	Use:   "process-orders",
	Short: "Issue orders or accept customer returns",
}

func SetupProcessCmd(orderService *order.Service) {
	processCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	processCmd.Flags().StringP(FlagAction, "a", "", "Expires")
	processCmd.Flags().StringP(FlagOrderIDs, "o", "", "Order ID")

	_ = processCmd.MarkFlagRequired(FlagOrderID)
	_ = processCmd.MarkFlagRequired(FlagUserID)
	_ = processCmd.MarkFlagRequired(FlagExpires)

	processCmd.Run = func(cmd *cobra.Command, args []string) {
		userID, err := GetFlagString(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		action, err := GetFlagString(cmd, FlagAction)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		orderIDs, err := GetFlagString(cmd, FlagOrderIDs)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		sliceOrderIDs := strings.Split(orderIDs, ",")

		for _, orderID := range sliceOrderIDs {
			if err = orderService.Process(userID, orderID, action); err != nil {
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
