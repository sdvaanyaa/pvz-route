package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
)

var acceptCmd = &cobra.Command{
	Use:   "accept-order",
	Short: "Accept the order from the courier",
}

func SetupAcceptCmd(orderService *services.OrderService) {
	acceptCmd.Flags().StringP(FlagOrderID, "o", "", "Order ID")
	acceptCmd.Flags().StringP(FlagUserID, "u", "", "User ID")
	acceptCmd.Flags().StringP(FlagExpires, "e", "", "Expires")

	_ = acceptCmd.MarkFlagRequired(FlagOrderID)
	_ = acceptCmd.MarkFlagRequired(FlagUserID)
	_ = acceptCmd.MarkFlagRequired(FlagExpires)

	acceptCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, err := GetFlag(cmd, FlagOrderID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		userID, err := GetFlag(cmd, FlagUserID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		expire, err := GetFlag(cmd, FlagExpires)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		if err = orderService.AcceptOrder(orderID, userID, expire); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		fmt.Printf("ORDER_ACCEPTED: %s\n", orderID)
	}
}

func init() {
	rootCmd.AddCommand(acceptCmd)
}
