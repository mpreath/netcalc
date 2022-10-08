package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mpreath/netcalc/pkg/network"
	"github.com/mpreath/netcalc/pkg/network/networknode"
	"log"
	"sync"

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
		networkAddress, err := utils.ParseAddress(args[0])
		if err != nil {
			log.Fatal(err)
		}
		networkMask, err := utils.ParseAddress(args[1])
		if err != nil {
			log.Fatal(err)
		}
		net, err := network.New(networkAddress, networkMask)
		if err != nil {
			log.Fatal(err)
		}
		// generate network from args
		node := networknode.New(net)

		if HOST_COUNT > 0 {
			err := SplitToHostCountThreaded(node, HOST_COUNT)
			if err != nil {
				log.Fatal(err)
			}

		} else if NET_COUNT > 0 {
			err = networknode.SplitToNetCount(node, NET_COUNT)
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
			printNetworkTree(node)
		}

	},
}

func SplitToHostCountThreaded(node *networknode.NetworkNode, host_count int) error {
	wg := new(sync.WaitGroup)

	valid, err := networknode.ValidForHostCount(node.Network, host_count)
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
			valid, err := networknode.ValidForHostCount(node.Subnets[0].Network, host_count)
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
					valid, err := networknode.ValidForHostCount(node.Subnets[0].Subnets[0].Network, host_count)
					if err != nil {
						log.Fatal(err)
					}

					if valid { // these subnets are valid
						return nil
					} else {
						wg.Add(4)
						go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[0], host_count)
						go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[1], host_count)
						go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[0], host_count)
						go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[1], host_count)
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

func SplitToHostCountWrapper(wg *sync.WaitGroup, node *networknode.NetworkNode, host_count int) {
	defer wg.Done()
	err := networknode.SplitToHostCount(node, host_count)
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	subnetCmd.Flags().IntVar(&HOST_COUNT, "hosts", 0, "Specifies the number of hosts to include each subnet.")
	subnetCmd.Flags().IntVar(&NET_COUNT, "networks", 0, "Specifies the number of subnets to create.")
	rootCmd.AddCommand(subnetCmd)
}

func printNetworkTree(node *networknode.NetworkNode, opts ...int) {
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
