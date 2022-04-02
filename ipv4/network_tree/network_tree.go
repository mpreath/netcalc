package network_tree

import (
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/network"
)

type NetworkNode struct {
	Network *network.Network `json:"network,omitempty"`
	Parent  *NetworkNode     `json:"-"`
	Subnets []*NetworkNode   `json:"subnets,omitempty"`
}

func (node *NetworkNode) Split() error {
	bc := ipv4.GetBitsInMask(node.Network.Mask) + 1
	if bc < 31 {
		new_mask, err := ipv4.GetMaskFromBits(bc)
		if err != nil {
			return err
		}
		// left will contain the lower value
		left_network, err := network.GenerateNetworkFromBits(node.Network.Address, new_mask)
		if err != nil {
			return err
		}

		// node.Subnets[0] = &NetworkNode{Parent: node, Network: left_network}

		node.Subnets = append(node.Subnets, &NetworkNode{Parent: node, Network: left_network})

		// right will contain the larger value
		right_network, err := network.GenerateNetworkFromBits(left_network.BroadcastAddress+1, new_mask)
		if err != nil {
			return err
		}

		// node.Subnets[0] = &NetworkNode{Parent: node, Network: right_network}

		node.Subnets = append(node.Subnets, &NetworkNode{Parent: node, Network: right_network})

		node.Network.Hosts = nil
	} else {
		return fmt.Errorf("Network doesn't support being split.\n")
	}

	return nil
}
