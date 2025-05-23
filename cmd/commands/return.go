package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var returnCmd = &cobra.Command{
	Use:   "return-order",
	Short: "Return the order to the courier",
}

func SetupReturnCmd(orderService *order.Service) {
	returnCmd.Flags().StringP(FlagOrderID, "o", "", "Order ID")

	_ = returnCmd.MarkFlagRequired(FlagOrderID)

	returnCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, err := GetFlagString(cmd, FlagOrderID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		if err = orderService.Return(orderID); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		fmt.Printf("ORDER_RETURNED: %s\n", orderID)
	}
}

func init() {
	rootCmd.AddCommand(returnCmd)
}
