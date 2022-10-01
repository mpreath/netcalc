package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mpreath/netcalc/pkg/network"
	"github.com/mpreath/netcalc/pkg/utils"
	"github.com/spf13/cobra"
)

var HOST_COUNT int
var NET_COUNT int

var subnetCmd = &cobra.Command{
	Use:   "subnet [--hosts <hosts> | --networks <networks>] <ip_address> <subnet_mask>",
	Short: "Given a network break it into smaller networks",
	Long: `
This command subnets a network based on host count and network count parameters.
Usage: netcalc subnet [--hosts <num of hosts>|--nets <num of networks>] <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		baseNetwork, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		// contains a list of summarized networks
		var summarizedNetworks []*network.Network

		if HOST_COUNT > 0 {
			summarizedNetworks, err = network.SplitToHostCountv2(baseNetwork, HOST_COUNT)
			if err != nil {
				log.Fatal(err)
			}

		} else if NET_COUNT > 0 {
			summarizedNetworks, err = network.SplitToNetCountv2(baseNetwork, NET_COUNT)
			if err != nil {
				log.Fatal(err)
			}
		}

		if JSON_FLAG {
			// json output
			s, _ := json.MarshalIndent(summarizedNetworks, "", "  ")
			fmt.Println(string(s))
		} else {
			// std output
			printNetworkSlice(summarizedNetworks)
		}

	},
}

func init() {
	subnetCmd.Flags().IntVar(&HOST_COUNT, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	subnetCmd.Flags().IntVar(&NET_COUNT, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(subnetCmd)
}

func printNetworkSlice(summary []*network.Network) {

	if VERBOSE_FLAG {
		fmt.Printf("Summarized into %d network(s)\n", len(summary))
	}
	for _, net := range summary {
		fmt.Printf("%s\t%s\n", utils.Itodd(net.Address), utils.Itodd(net.Mask))
	}
}

func printNetworkTree(node *network.NetworkNode, opts ...int) {
	var depth int

	if len(opts) == 0 {
		depth = 0
	} else {
		depth = opts[0]
	}

	if VERBOSE_FLAG {
		if depth == 0 {
			fmt.Printf("* = assigned network\n")
			fmt.Printf("+ = useable network\n")
			fmt.Printf("[n] = # of useable hosts\n\n")
		}

		ip_address := utils.Itodd(node.Network.Address)
		num_of_bits := utils.GetBitsInMask(node.Network.Mask)

		for i := 0; i < depth; i++ {
			fmt.Printf(" |")
		}

		fmt.Printf("__%s/%d", ip_address, num_of_bits)
		if node.Utilized && len(node.Subnets) == 0 {
			fmt.Printf("[%d]*", node.Network.HostCount())
		} else if len(node.Subnets) == 0 {
			fmt.Printf("[%d]+", node.Network.HostCount())
		}
		fmt.Printf("\n")
	} else {
		if len(node.Subnets) == 0 {
			ip_address := utils.Itodd(node.Network.Address)
			mask := utils.Itodd(node.Network.Mask)
			fmt.Printf("%s\t%s\n", ip_address, mask)
		}
	}

	if len(node.Subnets) > 0 {
		printNetworkTree(node.Subnets[0], depth+1)
		printNetworkTree(node.Subnets[1], depth+1)
	}

}
