package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/network"
	"github.com/mpreath/netcalc/utils"
	"github.com/spf13/cobra"
)

var HOST_COUNT int
var NET_COUNT int

var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Given a network break it into smaller networks",
	Long: `
This command subnets a network based on host count and network count parameters.
Usage: netcalc info <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		net, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		// generate network from args
		node := network.NetworkNode{
			Network: net,
		}

		if HOST_COUNT > 0 {
			err = network.SplitToHostCount(&node, HOST_COUNT)
			if err != nil {
				fmt.Println(err)
				return
			}

		} else if NET_COUNT > 0 {
			err = network.SplitToNetCount(&node, NET_COUNT)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if JSON_FLAG {
			// json output
			s, _ := json.MarshalIndent(node, "", "  ")
			fmt.Println(string(s))
		} else {
			// std output
			printNetworkTree(&node)
		}

	},
}

func init() {
	subnetCmd.Flags().IntVar(&HOST_COUNT, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	subnetCmd.Flags().IntVar(&NET_COUNT, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(subnetCmd)
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
			// fmt.Printf("* = assigned network\n")
			fmt.Printf("+ = useable network\n")
			fmt.Printf("[n] = # of useable hosts\n\n")
		}

		ip_address := utils.Itodd(node.Network.Address)
		num_of_bits := utils.GetBitsInMask(node.Network.Mask)

		for i := 0; i < depth; i++ {
			fmt.Printf(" |")
		}

		fmt.Printf("__%s/%d", ip_address, num_of_bits)
		if len(node.Subnets) == 0 {
			fmt.Printf("+[%d]", len(node.Network.Hosts))
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
