package main

import (
	"encoding/json"
	"fmt"
	"github.com/mpreath/netcalc/pkg/netcalc"
	"log"
	"sync"

	"github.com/spf13/cobra"
)

var HostCount int
var NetCount int

var subnetCmd = &cobra.Command{
	Use:   "subnet [--hosts <hosts> | --networks <networks>] <ip_address> <subnet_mask>",
	Short: "Given a network break it into smaller networks",
	Long: `
This command subnets a network based on host count and network count parameters.
Usage: netcalc subnet [--hosts <num of hosts>|--nets <num of networks>] <ip_address> <subnet_mask>.`,
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
		net, err := netcalc.NewNetwork(networkAddress, networkMask)
		if err != nil {
			log.Fatal(err)
		}
		// generate network from args
		node := netcalc.NewNetworkNode(net)

		if HostCount > 0 {
			err := SplitToHostCountThreaded(node, HostCount)
			if err != nil {
				log.Fatal(err)
			}

		} else if NetCount > 0 {
			err = netcalc.SplitToNetCount(node, NetCount)
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
			printNetworkTree(node)
		}

	},
}

func SplitToHostCountThreaded(node *netcalc.NetworkNode, hostCount int) error {
	wg := new(sync.WaitGroup)

	valid, err := netcalc.ValidForHostCount(node.Network, hostCount)
	if err != nil {
		log.Fatal(err)
	}

	if valid { // base network is already the best option
		return nil
	} else { // we can subnet another level
		err = node.Split() // create two subnets
		if err != nil {
			log.Fatal(err)
		}
		if len(node.Subnets) > 0 {
			valid, err := netcalc.ValidForHostCount(node.Subnets[0].Network, hostCount)
			if err != nil {
				log.Fatal(err)
			}

			if valid { // these subnets are valid
				return nil
			} else {
				err = node.Subnets[0].Split()
				if err != nil {
					log.Fatal(err)
				}

				err = node.Subnets[1].Split()
				if err != nil {
					log.Fatal(err)
				}

				if len(node.Subnets[0].Subnets) > 0 && len(node.Subnets[1].Subnets) > 0 {
					valid, err := netcalc.ValidForHostCount(node.Subnets[0].Subnets[0].Network, hostCount)
					if err != nil {
						log.Fatal(err)
					}

					if valid { // these subnets are valid
						return nil
					} else {
						wg.Add(4)
						go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[0], hostCount)
						go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[1], hostCount)
						go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[0], hostCount)
						go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[1], hostCount)
					}

				}
			}
		} else {
			return nil
		}
	}

	wg.Wait()

	return nil
}

func SplitToHostCountWrapper(wg *sync.WaitGroup, node *netcalc.NetworkNode, hostCount int) {
	defer wg.Done()
	err := netcalc.SplitToHostCount(node, hostCount)
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	subnetCmd.Flags().IntVar(&HostCount, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	subnetCmd.Flags().IntVar(&NetCount, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(subnetCmd)
}

func printNetworkTree(node *netcalc.NetworkNode, opts ...int) {
	var depth int

	if len(opts) == 0 {
		depth = 0
	} else {
		depth = opts[0]
	}

	if VerboseFlag {
		if depth == 0 {
			fmt.Printf("* = assigned network\n")
			fmt.Printf("+ = useable network\n")
			fmt.Printf("[n] = # of useable hosts\n\n")
		}

		ipAddress := netcalc.ExportAddress(node.Network.Address)
		numOfBits := netcalc.GetBitsInMask(node.Network.Mask)

		for i := 0; i < depth; i++ {
			fmt.Printf(" |")
		}

		fmt.Printf("__%s/%d", ipAddress, numOfBits)
		if node.Utilized && len(node.Subnets) == 0 {
			fmt.Printf("[%d]*", node.Network.HostCount())
		} else if len(node.Subnets) == 0 {
			fmt.Printf("[%d]+", node.Network.HostCount())
		}
		fmt.Printf("\n")
	} else {
		if len(node.Subnets) == 0 {
			ipAddress := netcalc.ExportAddress(node.Network.Address)
			mask := netcalc.ExportAddress(node.Network.Mask)
			fmt.Printf("%s\t%s\n", ipAddress, mask)
		}
	}

	if len(node.Subnets) > 0 {
		printNetworkTree(node.Subnets[0], depth+1)
		printNetworkTree(node.Subnets[1], depth+1)
	}

}
