package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

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
		net, err := network.GenerateNetwork(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		// generate network from args
		node := network.NetworkNode{
			Network: net,
		}

		// cpuFile, err := os.Create("tmp/cpuProfile8xGo.pprof")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// pprof.StartCPUProfile(cpuFile)
		// defer pprof.StopCPUProfile()

		if HOST_COUNT > 0 {
			//err = network.SplitToHostCount(&node, HOST_COUNT)
			SplitToHostCountThreaded(&node, HOST_COUNT)
			if err != nil {
				log.Fatal(err)
			}

		} else if NET_COUNT > 0 {
			err = network.SplitToNetCount(&node, NET_COUNT)
			if err != nil {
				log.Fatal(err)
			}
		}

		// pprof.StopCPUProfile()

		if JSON_FLAG {
			// json output
			s, _ := json.MarshalIndent(node, "", "  ")
			fmt.Println(string(s))
		} else {
			// std output
			// printNetworkTree(&node)
		}

	},
}

func SplitToHostCountThreaded(node *network.NetworkNode, host_count int) error {

	if node.Network.HostCount() >= 30 {
		return nil
	}

	wg := new(sync.WaitGroup)
	node.Split()
	if len(node.Subnets) > 0 {
		// wg.Add(2)
		// go SplitToHostCountWrapper(wg, node.Subnets[0], host_count)
		// go SplitToHostCountWrapper(wg, node.Subnets[1], host_count)
		node.Subnets[0].Split()
		node.Subnets[1].Split()

		if len(node.Subnets[0].Subnets) > 0 && len(node.Subnets[1].Subnets) > 0 {

			wg.Add(4)
			go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[0], host_count)
			go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[1], host_count)
			go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[0], host_count)
			go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[1], host_count)

			// node.Subnets[0].Subnets[0].Split()
			// node.Subnets[0].Subnets[1].Split()
			// node.Subnets[1].Subnets[0].Split()
			// node.Subnets[1].Subnets[1].Split()

			// if len(node.Subnets[0].Subnets[0].Subnets) > 0 && len(node.Subnets[0].Subnets[1].Subnets) > 0 && len(node.Subnets[0].Subnets[0].Subnets) > 0 && len(node.Subnets[0].Subnets[0].Subnets) > 0 {
			// 	wg.Add(8)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[0].Subnets[0], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[0].Subnets[1], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[1].Subnets[0], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[0].Subnets[1].Subnets[1], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[0].Subnets[0], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[0].Subnets[1], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[1].Subnets[0], host_count)
			// 	go SplitToHostCountWrapper(wg, node.Subnets[1].Subnets[1].Subnets[1], host_count)
			// }
		}
	}
	wg.Wait()

	return nil
}

func SplitToHostCountWrapper(wg *sync.WaitGroup, node *network.NetworkNode, host_count int) {
	defer wg.Done()
	network.SplitToHostCount(node, host_count)

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
