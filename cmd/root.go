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
		// TODO: Get arguments for address and mask
		n, err := network.GenerateNetwork("192.168.1.1", "255.255.255.0")
		if err != nil {
			fmt.Print(err)
			return
		}
		// TODO: Add flags for verbosity and JSON
		printNetworkInformation(n, false)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func printNetworkInformation(n *network.Network, verbose bool) {
	if n != nil {
		n_dd_address := ipv4.Itodd(n.Address)
		n_dd_mask := ipv4.Itodd(n.Mask)
		n_dd_bcast := ipv4.Itodd(n.BroadcastAddress)

		fmt.Printf("Network:\t%s\n", n_dd_address)
		fmt.Printf("Mask:\t\t%s\n", n_dd_mask)
		fmt.Printf("Broadcast:\t%s\n", n_dd_bcast)
		fmt.Printf("Usable Hosts:\t%d\n", len(n.Hosts))

		if verbose {
			// TODO: Print out network hosts
		}
	}
}

func printNetworkInformationJSON(network *network.Network) {
	// TODO: Create appropriate marshalling for Network and Host structs
}
