package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services"
)

var returnCmd = &cobra.Command{
	Use:   "return-order",
	Short: "Return the order to the courier",
}

func SetupReturnCmd(orderService *services.OrderService) {
	returnCmd.Flags().StringP(FlagOrderID, "o", "", "Order ID")

	_ = returnCmd.MarkFlagRequired(FlagOrderID)

	returnCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, err := GetFlag(cmd, FlagOrderID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		if err = orderService.ReturnOrder(orderID); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}

		fmt.Printf("ORDER_RETURNED: %s\n", orderID)
	}
}

func init() {
	rootCmd.AddCommand(returnCmd)
}
