package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/network"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := network.GenerateNetwork(args[0], args[1])
		s, _ := json.MarshalIndent(n, "", "  ")
		fmt.Printf("%s\n", s)
	},
}
