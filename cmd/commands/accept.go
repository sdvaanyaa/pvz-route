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

	acceptCmd.RunE = func(cmd *cobra.Command, args []string) error {
		orderID, err := GetFlag(cmd, FlagOrderID)
		if err != nil {
			return err
		}

		userID, err := GetFlag(cmd, FlagUserID)
		if err != nil {
			return err
		}

		expire, err := GetFlag(cmd, FlagExpires)
		if err != nil {
			return err
		}

		if err := orderService.AcceptOrder(orderID, userID, expire); err != nil {
			return err
		}

		fmt.Printf("Order %s accepted\n", orderID)

		return nil
	}
}

func init() {
	rootCmd.AddCommand(acceptCmd)
}
