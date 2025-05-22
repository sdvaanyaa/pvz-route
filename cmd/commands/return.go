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

	returnCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		orderID, err := GetFlag(cmd, FlagOrderID)
		if err != nil {
			return err
		}

		if err = orderService.ReturnOrder(orderID); err != nil {
			return err
		}

		fmt.Printf("ORDER_RETURNED: %s\n", orderID)

		return nil
	}
}

func init() {
	rootCmd.AddCommand(returnCmd)
}
