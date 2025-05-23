package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var acceptCmd = &cobra.Command{
	Use:   "accept-order",
	Short: "Accept the order from the courier",
}

func SetupAcceptCmd(orderService *order.Service) {
	acceptCmd.Flags().StringP(FlagOrderID, "o", "", "Order ID")
	acceptCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	acceptCmd.Flags().StringP(FlagExpires, "e", "", "Expires")

	_ = acceptCmd.MarkFlagRequired(FlagOrderID)
	_ = acceptCmd.MarkFlagRequired(FlagUserID)
	_ = acceptCmd.MarkFlagRequired(FlagExpires)

	acceptCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, err := GetFlagString(cmd, FlagOrderID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		userID, err := GetFlagString(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		expire, err := GetFlagString(cmd, FlagExpires)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		if err = orderService.Accept(orderID, userID, expire); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		fmt.Printf("ORDER_ACCEPTED: %s\n", orderID)
	}
}

func init() {
	rootCmd.AddCommand(acceptCmd)
}
