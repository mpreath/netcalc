package netcalc

import (
	"fmt"
	"sort"
)

// SummarizeNetworks takes an array of Network objects and returns a summarized Network.
// Summarization occurs by sorting the array and then finding a common bit boundary that
// a common mask and common network bits can be calculated.
func SummarizeNetworks(networks []*Network) (*Network, error) {

	if len(networks) == 0 {
		return nil, fmt.Errorf("SummarizeNetwork: no networks to summarize")
	} else if len(networks) == 1 {
		return networks[0], nil
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Address < networks[j].Address
	})

	commonBits := networks[0].Address
	var commonMask uint32

	for idx := 1; idx < len(networks); idx++ {
		commonBits = commonBits & networks[idx].Address
		commonMask = GetCommonBitMask(commonBits, networks[idx].Address)
		commonBits = GetNetworkAddress(commonBits, commonMask)
	}

	return &Network{Address: commonBits, Mask: commonMask}, nil

}
