package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	FlagOrderID  = "order-id"
	FlagUserID   = "user-id"
	FlagExpires  = "expires"
	FlagLimit    = "limit"
	FlagAction   = "action"
	FlagOrderIDs = "order-ids"
	FlagInPVZ    = "in-pvz"
	FlagLast     = "last"
	FlagPage     = "page"
	FlagFile     = "file"
)

func GetFlag(cmd *cobra.Command, name string) (string, error) {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", fmt.Errorf("cannot read flag --%s: %w", name, err)
	}
	return val, nil
}
