package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var returnCmd = &cobra.Command{
	Use:   "return-order",
	Short: "Return the order to the courier",
}

func setupReturnCmd(orderSvc order.Service) {
	returnCmd.Flags().StringP(FlagOrderID, ShortOrderID, "", "Order ID")

	_ = returnCmd.MarkFlagRequired(FlagOrderID)

	returnCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, err := getFlagString(cmd, FlagOrderID)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		if err = orderSvc.Return(orderID); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		fmt.Printf("ORDER_RETURNED: %s\n", orderID)
	}
}

func init() {
	rootCmd.AddCommand(returnCmd)
}
