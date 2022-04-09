package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mpreath/netcalc/network"
	"github.com/mpreath/netcalc/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
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

		// TODO: Determine if sorting helps performance of the overall algorithm
		// sort the slice based on network address (uint32)
		sort.Slice(networks, func(i, j int) bool {
			return networks[i].Address < networks[j].Address
		})

		// Summarization Method A

		// tracks whether summarization has occured within inner loop
		// if changes made we should loop through the networks slice again
		// to see if more summarization is necessary
		changes_made := true

		for changes_made {
			// outer loop
			// once inside the outer loop assume no changes have been made
			// if we sumamrize in the inner loop we will update
			changes_made = false
			for oidx := 0; oidx < len(networks); oidx++ {
				//fmt.Printf("[%d] ", oidx)
				if networks[oidx] == nil {
					// if the network at the outer index is nil
					// move onto the next
					//fmt.Printf("(nil) \n")
					continue
				}

				//fmt.Printf("\n")
				// inner loop
				for iidx := 0; iidx < len(networks); iidx++ {
					//fmt.Printf("[%d] [%d] ", oidx, iidx)
					// we compare the network at the outer index (oidx)
					// to the network at the inner index (iidx) to determine
					// if summarization is possible
					if networks[iidx] == nil || iidx == oidx {
						// if the network at the inner index is nil
						// move onto the next
						if networks[iidx] == nil {
							//fmt.Printf("(nil) \n")
						}
						if iidx == oidx {
							//fmt.Printf("(same) \n")
						}
						continue
					}

					if networks[iidx].Mask != networks[oidx].Mask {
						//fmt.Printf("(mask mismatch) \n")
						continue
					}

					// determine if these two networks can be summarized
					new_mask := GetCommonBitMask(networks[oidx].Address, networks[iidx].Address)

					// if these two networks can be summarized new_mask will return
					// a new mask and if not will return 0
					if new_mask != 0 {
						//fmt.Printf("(summarized) [%s][%s] [%s][%s]\n", utils.Itodd(networks[oidx].Address), utils.Itodd(networks[oidx].Mask), utils.Itodd(networks[iidx].Address), utils.Itodd(networks[iidx].Mask))
						// update the base/outside network to be a summary
						// of the two networks (outside/inside) by reducing the
						// mask by a bit
						networks[oidx] = &network.Network{
							Address: utils.GetNetworkAddress(networks[oidx].Address, new_mask),
							Mask:    new_mask,
						}
						// set the inside network to nil
						// because it is now summarized into the outside network
						networks[iidx] = nil

						// set changes_made to true
						changes_made = true

						// since we made a summarization, lets move to the next outside network

						break
					} else {
						//fmt.Printf("(can't summarize) [%s][%s] [%s][%s]\n", utils.Itodd(networks[oidx].Address), utils.Itodd(networks[oidx].Mask), utils.Itodd(networks[iidx].Address), utils.Itodd(networks[iidx].Mask))
					}

					//fmt.Printf("\n")
				}
			}
		}

		// test summarization
		// changes_made := true
		// outer_index := 1
		// fmt.Printf("### Summarizing ###\n")
		// for changes_made {

		// 	changes_made = false
		// 	fmt.Printf(">[%d]\n", outer_index)

		// 	for idx := 0; idx < len(networks); idx++ {
		// 		fmt.Printf("\t[%d] [%d] [%d]\n", idx, outer_index)
		// 		if networks[idx] == nil {
		// 			fmt.Printf("\tskipped\n")
		// 			continue
		// 		}

		// 		if idx+idx_offset < len(networks) && networks[idx+idx_offset] != nil {
		// 			first_value := networks[idx].Address
		// 			second_value := networks[idx+idx_offset].Address
		// 			new_mask := GetCommonBitMask(first_value, second_value)

		// 			if new_mask != 0 {
		// 				networks[idx] = &network.Network{
		// 					Address: first_value,
		// 					Mask:    new_mask,
		// 				}
		// 				networks[idx+idx_offset] = nil
		// 				idx = idx + idx_offset
		// 				changes_made = true
		// 			} else {
		// 				// we hit a bit boundry change (sorted)
		// 			}
		// 		} else {
		// 			//fmt.Printf("[%d] %s\n", idx, utils.Itodd(networks[idx].Address))
		// 		}

		// 	}
		// 	idx_offset = idx_offset + idx_offset - 1
		// 	outer_index++

		// }
		// fmt.Printf("### End ###\n")

		// print
		for idx, net := range networks {
			if net != nil {
				fmt.Printf("[%d]\t%s\t%s\n", idx, utils.Itodd(net.Address), utils.Itodd(net.Mask))
			}
		}
	},
}

func SummarizeNetworkSlice(networks []network.Network) []network.Network {

	return nil
}

func GetCommonBitMask(n1 uint32, n2 uint32) uint32 {
	common_bits := n1 ^ n2

	//fmt.Printf("\t%s <-> %s\n", utils.Itodd(n1), utils.Itodd(n2))

	if CheckNumberPowerOfTwo(common_bits) {
		// if its a power of two it means only one bit is set
		// lets determine where that it
		for idx := 0; idx < 32; idx++ {
			//fmt.Printf("(%d,%d)\n", idx, common_bits)
			if common_bits == 0 {
				// found it
				// fmt.Printf("found it\n")
				new_mask, _ := utils.GetMaskFromBits(32 - idx)
				return new_mask
			}
			common_bits = common_bits >> 1
		}

		return 0
	} else {
		// fmt.Printf("not a power of two\n")
		return 0
	}
}

func CheckNumberPowerOfTwo(n uint32) bool {
	val := n & (n - 1)
	if val == 0 {
		return true
	} else {
		return false
	}
}
