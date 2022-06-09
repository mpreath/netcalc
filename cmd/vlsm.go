package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mpreath/netcalc/network"
	"github.com/spf13/cobra"
)

var vlsmCmd = &cobra.Command{
	Use:   "vlsm",
	Short: "Given a network and comma-separated list of subnet lengths break it into smaller networks",
	Long: `
This command subnets a network based on a comma-separated list of subnet lengths.
Usage: netcalc vlsm <vlsm list> <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		net, err := network.GenerateNetwork(args[1], args[2])
		if err != nil {
			log.Fatal(err)
		}
		// generate network from args
		node := network.NetworkNode{
			Network: net,
		}

		vlsm_list := strings.Split(args[0], ",")

		for _, vlsm := range vlsm_list {
			host_count, _ := strconv.ParseInt(vlsm, 10, 32)
			fmt.Printf("%d\n", host_count)
			err = network.SplitToHostCount(&node, int(host_count))

			if err != nil {
				log.Fatal(err)
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
	// subnetCmd.Flags().IntVar(&HOST_COUNT, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	// subnetCmd.Flags().IntVar(&NET_COUNT, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(vlsmCmd)
}
