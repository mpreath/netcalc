package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpreath/netcalc/pkg/netcalc"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info <ip_address> <subnet_mask>",
	Short: "Displays information about a network",
	Long: `
This command displays information about a network.
Usage: netcalc info <ip_address> <subnet_mask>.`,

	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		networkAddress, err := netcalc.ParseAddress(args[0])
		if err != nil {
			log.Fatal(err)
		}
		networkMask, err := netcalc.ParseAddress(args[1])
		if err != nil {
			log.Fatal(err)
		}
		n, err := netcalc.NewNetwork(networkAddress, networkMask)
		if err != nil {
			log.Fatal(err)
		}

		if JsonFlag {
			err := printNetworkInformationJSON(n)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			printNetworkInformation(n)
		}
	},
}

func printNetworkInformation(n *netcalc.Network) {
	if n != nil {

		fmt.Printf("Network:\t%s\n", netcalc.ExportAddress(n.Address))
		fmt.Printf("Mask:\t\t%s (/%d)\n", netcalc.ExportAddress(n.Mask), netcalc.GetBitsInMask(n.Mask))
		fmt.Printf("Bcast:\t\t%s\n", netcalc.ExportAddress(n.BroadcastAddress()))

		if VerboseFlag {
			fmt.Printf("\n")

			for _, host := range n.Hosts() {
				fmt.Printf("%s\t%s\n", netcalc.ExportAddress(host.Address), netcalc.ExportAddress(host.Mask))
			}
		}

	}
}

func printNetworkInformationJSON(network *netcalc.Network) error {
	s, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)
	return nil
}
