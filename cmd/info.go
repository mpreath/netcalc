package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/network"
	"github.com/mpreath/netcalc/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about a network",
	Long: `
This command displays information about an utils network.
Usage: netcalc info <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			fmt.Print(err)
			return
		}

		if JSON_FLAG {
			printNetworkInformationJSON(n)
		} else {
			printNetworkInformation(n)
		}
	},
}

func printNetworkInformation(n *network.Network) {
	if n != nil {

		fmt.Printf("Network:\t%s\n", utils.Itodd(n.Address))
		fmt.Printf("Mask:\t\t%s\n", utils.Itodd(n.Mask))
		fmt.Printf("Broadcast:\t%s\n", utils.Itodd(n.BroadcastAddress))
		fmt.Printf("Usable Hosts:\t%d\n", len(n.Hosts))

		if VERBOSE_FLAG {
			for _, host := range n.Hosts {
				fmt.Printf("%s\t%s\n", utils.Itodd(host.Address), utils.Itodd(host.Mask))
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