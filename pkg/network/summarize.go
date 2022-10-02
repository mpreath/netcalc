package network

import (
	"fmt"
	"sort"

	"github.com/mpreath/netcalc/pkg/utils"
)

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

	for oidx := 1; oidx < len(networks); oidx++ {
		commonBits = commonBits & networks[oidx].Address
		commonMask = utils.GetCommonBitMask(commonBits, networks[oidx].Address)
		commonBits = utils.GetNetworkAddress(commonBits, commonMask)
	}

	return &Network{Address: commonBits, Mask: commonMask}, nil

}
