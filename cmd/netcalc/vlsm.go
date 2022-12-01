package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpreath/netcalc/pkg/netcalc"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var vlsmCmd = &cobra.Command{
	Use:   "vlsm <host_counts_list> <ip_address> <subnet_mask>",
	Short: "Given a network and comma-separated list of subnet host counts break it into smaller networks",
	Long: `
This command subnets a network based on a comma-separated list of subnet host counts.
Usage: netcalc vlsm <host_counts_list> <ip_address> <subnet_mask>.`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		networkAddress, err := netcalc.ParseAddress(args[1])
		if err != nil {
			log.Fatal(err)
		}
		networkMask, err := netcalc.ParseAddress(args[2])
		if err != nil {
			log.Fatal(err)
		}
		net, err := netcalc.NewNetwork(networkAddress, networkMask)
		if err != nil {
			log.Fatal(err)
		}
		// generate network from args
		node := netcalc.NetworkNode{
			Network: net,
		}

		vlsmArgs := strings.Split(args[0], ",")
		var vlsmList = make([]int, len(vlsmArgs))
		for idx, val := range vlsmArgs {
			vlsmList[idx], err = strconv.Atoi(val)
			if err != nil {
				log.Fatal(err)
			}
		}
		sort.Slice(vlsmList, func(i, j int) bool {
			return vlsmList[i] < vlsmList[j]
		})

		for _, vlsm := range vlsmList {
			err = netcalc.SplitToVlsmCount(&node, vlsm)

			if err != nil {
				log.Fatal(err)
			}
		}

		if JsonFlag {
			// json output
			s, _ := json.MarshalIndent(node, "", "  ")
			fmt.Println(string(s))
		} else {
			// std output
			printNetworkTree(&node)
			s, _ := json.MarshalIndent(node.FlattenUtilized(), "", "  ")
			fmt.Println(string(s))
		}

	},
}

func init() {
	rootCmd.AddCommand(vlsmCmd)
}
