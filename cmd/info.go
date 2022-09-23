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
	Use:   "info",
	Short: "Displays information about a network",
	Long: `
This command displays information about an utils network.
Usage: netcalc info <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		n, err := network.GenerateNetwork(args[0], args[1])
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

		if VERBOSE_FLAG {
			fmt.Printf("* = network address\n")
			fmt.Printf("+ = broadcast address\n\n")
			fmt.Printf("%s\t%s *\n", utils.Itodd(n.Address), utils.Itodd(n.Mask))
		}

		for _, host := range network.GetHosts(n) {
			fmt.Printf("%s\t%s\n", utils.Itodd(host.Address), utils.Itodd(host.Mask))
		}

		if VERBOSE_FLAG {
			fmt.Printf("%s\t%s +\n", utils.Itodd(n.BroadcastAddress), utils.Itodd(n.Mask))
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
