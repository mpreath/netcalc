package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/ipv4/network"
	"github.com/mpreath/netcalc/ipv4/network_tree"
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
		network, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		// generate network from args
		node := network_tree.NetworkNode{
			Parent:  nil,
			Network: network,
		}
		err = node.Split()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = node.Subnets[0].Split()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = node.Subnets[1].Split()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = node.Subnets[0].Subnets[0].Split()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = node.Subnets[0].Subnets[0].Subnets[0].Split()
		if err != nil {
			fmt.Println(err)
			return
		}

		//_ = json.NewEncoder(os.Stdout).Encode(node)
		s, _ := json.MarshalIndent(node, "", "  ")
		fmt.Println(string(s))

	},
}

func init() {
	subnetCmd.Flags().IntVar(&HOST_COUNT, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	subnetCmd.Flags().IntVar(&NET_COUNT, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(subnetCmd)
}
