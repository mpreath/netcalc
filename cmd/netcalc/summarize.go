package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mpreath/netcalc/pkg/netcalc"
	"github.com/mpreath/netcalc/pkg/netcalc/networknode"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "summarizes the networks provided to stdin",
	Run: func(cmd *cobra.Command, args []string) {
		var networks []*netcalc.Network

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := strings.Split(scanner.Text(), "\t")

			networkAddress, err := netcalc.ParseAddress(input[0])
			if err != nil {
				log.Fatal(err)
			}
			networkMask, err := netcalc.ParseAddress(input[1])
			if err != nil {
				log.Fatal(err)
			}
			net, err := netcalc.NewNetwork(networkAddress, networkMask)
			if err != nil {
				log.Fatal(err)
			}

			networks = append(networks, net)
		}

		networkSummary, err := netcalc.SummarizeNetworks(networks)

		if err != nil {
			log.Fatal(err)
		}

		node := networknode.NetworkNode{
			Network: networkSummary,
		}

		if JsonFlag {
			// json output
			s, _ := json.MarshalIndent(node, "", "  ")
			fmt.Println(string(s))
		} else {
			// std output
			printNetworkTree(&node)
		}
	},
}
