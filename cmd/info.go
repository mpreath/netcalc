package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mpreath/netcalc/pkg/network"
	"github.com/mpreath/netcalc/pkg/utils"
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
		networkAddress, err := utils.ParseAddress(args[0])
		if err != nil {
			log.Fatal(err)
		}
		networkMask, err := utils.ParseAddress(args[1])
		if err != nil {
			log.Fatal(err)
		}
		n, err := network.New(networkAddress, networkMask)
		if err != nil {
			log.Fatal(err)
		}

		if JSON_FLAG {
			err := printNetworkInformationJSON(n)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			printNetworkInformation(n)
		}
	},
}

func printNetworkInformation(n *network.Network) {
	if n != nil {

		fmt.Printf("Network:\t%s\n", utils.Itodd(n.Address))
		fmt.Printf("Mask:\t\t%s (/%d)\n", utils.Itodd(n.Mask), utils.GetBitsInMask(n.Mask))
		fmt.Printf("Bcast:\t\t%s\n", utils.Itodd(n.BroadcastAddress()))

		if VERBOSE_FLAG {
			fmt.Printf("\n")

			for _, host := range n.Hosts() {
				fmt.Printf("%s\t%s\n", utils.Itodd(host.Address), utils.Itodd(host.Mask))
			}
		}

	}
}

func printNetworkInformationJSON(network *network.Network) error {
	s, err := json.MarshalIndent(network, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)
	return nil
}
