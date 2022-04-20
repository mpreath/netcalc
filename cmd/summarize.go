package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mpreath/netcalc/network"
	"github.com/mpreath/netcalc/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "summarizes the networks provided to stdin",
	Run: func(cmd *cobra.Command, args []string) {
		var networks []*network.Network

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := strings.Split(scanner.Text(), "\t")

			address, err := utils.Ddtoi(input[0])
			if err != nil {
				log.Fatal(err)
			}

			mask, err := utils.Ddtoi(input[1])
			if err != nil {
				log.Fatal(err)
			}

			networks = append(networks, &network.Network{
				Address: address,
				Mask:    mask,
			})
		}

		// print
		for idx, net := range network.SummarizeNetworks(networks) {
			if net != nil {
				fmt.Printf("[%d]\t%s\t%s\n", idx, utils.Itodd(net.Address), utils.Itodd(net.Mask))
			}
		}
	},
}
