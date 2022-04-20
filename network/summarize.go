package network

import (
	"sort"

	"github.com/mpreath/netcalc/utils"
)

func SummarizeNetworks(networks []*Network) []*Network {

	// Goal: Summarize the networks into as few networks as possible
	// Note: This should be O(n) and support multiple networks in the result
	// if they don't have a common parent

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
			if networks[oidx] == nil {
				// if the network at the outer index is nil
				// move onto the next
				continue
			}

			// inner loop
			for iidx := 0; iidx < len(networks); iidx++ {
				// we compare the network at the outer index (oidx)
				// to the network at the inner index (iidx) to determine
				// if summarization is possible
				if networks[iidx] == nil || iidx == oidx {
					// if the network at the inner index is nil
					// move onto the next
					continue
				}

				if networks[iidx].Mask != networks[oidx].Mask {
					continue
				}

				// determine if these two networks can be summarized
				new_mask := GetCommonBitMask(networks[oidx].Address, networks[iidx].Address)

				// if these two networks can be summarized new_mask will return
				// a new mask and if not will return 0
				if new_mask != 0 {
					// update the base/outside network to be a summary
					// of the two networks (outside/inside) by reducing the
					// mask by a bit
					networks[oidx] = &Network{
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
				}
			}
		}
	}

	return networks
}

func GetCommonBitMask(n1 uint32, n2 uint32) uint32 {
	common_bits := n1 ^ n2

	if CheckNumberPowerOfTwo(common_bits) {
		// if its a power of two it means only one bit is set
		// lets determine where that it
		for idx := 0; idx < 32; idx++ {
			if common_bits == 0 {
				// found it
				new_mask, _ := utils.GetMaskFromBits(32 - idx)
				return new_mask
			}
			common_bits = common_bits >> 1
		}

		return 0
	} else {
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
