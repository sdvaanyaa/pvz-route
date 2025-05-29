package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ozon.dev/sd_vaanyaa/homework/internal/services/order"
)

var acceptCmd = &cobra.Command{
	Use:   "accept-order",
	Short: "Accept the order from the courier",
}

func setupAcceptCmd(orderSvc order.Service) {
	acceptCmd.Flags().StringP(FlagOrderID, ShortOrderID, "", "Order ID")
	acceptCmd.Flags().StringP(FlagUserID, ShortUserID, "", "User ID")
	acceptCmd.Flags().StringP(FlagExpires, ShortExpires, "", "Expires")
	acceptCmd.Flags().Float64P(FlagWeight, ShortWeight, 0, "Weight")
	acceptCmd.Flags().Float64P(FlagPrice, ShortPrice, 0, "Price")
	acceptCmd.Flags().StringP(FlagPackageType, ShortPackageType, "", "Package type")

	_ = acceptCmd.MarkFlagRequired(FlagOrderID)
	_ = acceptCmd.MarkFlagRequired(FlagUserID)
	_ = acceptCmd.MarkFlagRequired(FlagExpires)
	_ = acceptCmd.MarkFlagRequired(FlagWeight)
	_ = acceptCmd.MarkFlagRequired(FlagPrice)

	acceptCmd.Run = func(cmd *cobra.Command, args []string) {
		orderID, userID, expire, packageType, weight, price, err := parseAcceptFlags(cmd)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		o, err := orderSvc.Accept(orderID, userID, expire, weight, price, packageType)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		fmt.Printf("ORDER_ACCEPTED: %s\n", o.ID)
		fmt.Printf("PACKAGE: %s\n", o.PackageType)
		fmt.Printf("TOTAL_PRICE: %.2f\n", o.Price)
	}
}

func init() {
	rootCmd.AddCommand(acceptCmd)
}

func parseAcceptFlags(cmd *cobra.Command) (string, string, string, string, float64, float64, error) {
	orderID, err := getFlagString(cmd, FlagOrderID)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	userID, err := getFlagString(cmd, FlagUserID)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	expire, err := getFlagString(cmd, FlagExpires)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	weight, err := getFlagFloat64(cmd, FlagWeight)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	price, err := getFlagFloat64(cmd, FlagPrice)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	packageType, err := getFlagString(cmd, FlagPackageType)
	if err != nil {
		return "", "", "", "", 0, 0, err
	}

	return orderID, userID, expire, packageType, weight, price, nil
}
