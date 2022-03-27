package cmd

import (
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/network"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netcalc",
	Short: "Netcalc is a IPv4/IPv6 network calculator",
	Run: func(cmd *cobra.Command, args []string) {
		// do stuff here
		n, err := network.GenerateNetwork("192.168.1.90", "255.255.255.252")
		if err != nil {
			fmt.Print(err)
			return
		}
		n_dd_address, _ := ipv4.Itodd(n.Address)
		n_dd_mask, _ := ipv4.Itodd(n.Mask)
		n_dd_bcast, _ := ipv4.Itodd(n.BroadcastAddress)
		fmt.Printf("Network: { Address: \"%s\", Mask: \"%s\", BroadcastAddress: \"%s\", Host Count: %d }\n", n_dd_address, n_dd_mask, n_dd_bcast, len(n.Hosts))
	},
}

func Execute() error {
	return rootCmd.Execute()
}
