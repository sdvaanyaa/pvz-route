package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagOrderID     = "order-id"
	FlagUserID      = "user-id"
	FlagExpires     = "expires"
	FlagLimit       = "limit"
	FlagAction      = "action"
	FlagOrderIDs    = "order-ids"
	FlagInPVZ       = "in-pvz"
	FlagLast        = "last"
	FlagPage        = "page"
	FlagFile        = "file"
	FlagPrice       = "price"
	FlagWeight      = "weight"
	FlagPackageType = "package"
)

const (
	ShortOrderID     = "o"
	ShortUserID      = "u"
	ShortExpires     = "e"
	ShortInPVZ       = "i"
	ShortLast        = "l"
	ShortPage        = "n"
	ShortLimit       = "m"
	ShortAction      = "a"
	ShortFile        = "f"
	ShortPrice       = "p"
	ShortWeight      = "w"
	ShortPackageType = "t"
)

var globalFlagSet = map[string]bool{
	"--" + FlagOrderID:  true,
	"-" + ShortOrderID:  true,
	"--" + FlagUserID:   true,
	"-" + ShortUserID:   true,
	"--" + FlagExpires:  true,
	"-" + ShortExpires:  true,
	"--" + FlagInPVZ:    true,
	"-" + ShortInPVZ:    true,
	"--" + FlagLast:     true,
	"-" + ShortLast:     true,
	"--" + FlagPage:     true,
	"-" + ShortPage:     true,
	"--" + FlagLimit:    true,
	"-" + ShortLimit:    true,
	"--" + FlagOrderIDs: true,
	"--" + FlagAction:   true,
	"-" + ShortAction:   true,
	"--" + FlagFile:     true,
	"-" + ShortFile:     true,
}

func validateFlagValue(val string) error {
	if globalFlagSet[val] {
		return fmt.Errorf("flag value cannot be another flag: %s", val)
	}
	return nil
}

func getFlagString(cmd *cobra.Command, name string) (string, error) {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", fmt.Errorf("cannot read flag --%s: %w", name, err)
	}

	if err = validateFlagValue(val); err != nil {
		return "", err
	}

	return val, nil
}

func getFlagBool(cmd *cobra.Command, name string) (bool, error) {
	val, err := cmd.Flags().GetBool(name)
	if err != nil {
		return false, fmt.Errorf("cannot read flag --%s: %w", name, err)
	}
	return val, nil
}

func getFlagInt(cmd *cobra.Command, name string) (int, error) {
	val, err := cmd.Flags().GetInt(name)
	if err != nil {
		return 0, fmt.Errorf("cannot read flag --%s: %w", name, err)
	}
	return val, nil
}

func getFlagFloat64(cmd *cobra.Command, name string) (float64, error) {
	val, err := cmd.Flags().GetFloat64(name)
	if err != nil {
		return 0, fmt.Errorf("cannot read flag --%s: %w", name, err)
	}
	return val, nil
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
	})
}
