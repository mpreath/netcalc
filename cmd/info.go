package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/network"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about a network",
	Long: `
This command displays information about an IPv4 network.
Usage: netcalc info <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			fmt.Print(err)
			return
		}

		if json_out {
			printNetworkInformationJSON(n)
		} else {
			printNetworkInformation(n)
		}
	},
}

func printNetworkInformation(n *network.Network) {
	if n != nil {

		fmt.Printf("Network:\t%s\n", ipv4.Itodd(n.Address))
		fmt.Printf("Mask:\t\t%s\n", ipv4.Itodd(n.Mask))
		fmt.Printf("Broadcast:\t%s\n", ipv4.Itodd(n.BroadcastAddress))
		fmt.Printf("Usable Hosts:\t%d\n", len(n.Hosts))

		if verbose {
			for _, host := range n.Hosts {
				fmt.Printf("%s\t%s\n", ipv4.Itodd(host.Address), ipv4.Itodd(host.Mask))
			}
		}
	}
}

func printNetworkInformationJSON(network *network.Network) {
	s, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("%s\n", s)
}
