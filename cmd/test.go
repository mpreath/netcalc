package cmd

import (
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	Run: func(cmd *cobra.Command, args []string) {
		mask := "128.0.0.0"
		uint_mask, _ := ipv4.Ddtoi(mask)
		bc := ipv4.GetBitsInMask(uint_mask)

		fmt.Printf("%s has %d bits\n", mask, bc)

		new_mask, _ := ipv4.GetMaskFromBits(bc)

		fmt.Printf("%s [%d]\n", ipv4.Itodd(new_mask), new_mask)
	},
}
