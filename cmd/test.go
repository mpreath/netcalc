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
		mask := "255.255.255.0"
		uint_mask, _ := ipv4.Ddtoi(mask)
		bc := ipv4.GetBitsInMask(uint_mask)

		fmt.Printf("%s has %d bits\n", mask, bc)

		bc = 25
		new_mask, err := ipv4.GetMaskFromBits(bc)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s [%d]\n", ipv4.Itodd(new_mask), new_mask)
	},
}
