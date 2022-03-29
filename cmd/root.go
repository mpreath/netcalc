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
		n, err := network.GenerateNetwork("192.168.1.1", "255.255.255.0")
		if err != nil {
			fmt.Print(err)
			return
		}
		printNetworkInformation(n)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func printNetworkInformation(n *network.Network) {
	n_dd_address, _ := ipv4.Itodd(n.Address)
	n_dd_mask, _ := ipv4.Itodd(n.Mask)
	n_dd_bcast, _ := ipv4.Itodd(n.BroadcastAddress)

	fmt.Printf("Network:\t%s\n", n_dd_address)
	fmt.Printf("Mask:\t\t%s\n", n_dd_mask)
	fmt.Printf("Broadcast:\t%s\n", n_dd_bcast)
	fmt.Printf("Usable Hosts:\t%d\n", len(n.Hosts))
}

func printNetworkInformationJSON(network *network.Network) {
}
