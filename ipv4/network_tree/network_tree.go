package network_tree

import (
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/network"
)

type NetworkNode struct {
	Network *network.Network `json:"network,omitempty"`
	Parent  *NetworkNode     `json:"-"`
	Left    *NetworkNode     `json:"subnet_1,omitempty"`
	Right   *NetworkNode     `json:"subnet_2,omitempty"`
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
		node.Left = &NetworkNode{
			Parent:  node,
			Left:    nil,
			Right:   nil,
			Network: left_network,
		}

		// right will contain the larger value
		right_network, err := network.GenerateNetworkFromBits(left_network.BroadcastAddress+1, new_mask)
		if err != nil {
			return err
		}
		node.Right = &NetworkNode{
			Parent:  node,
			Left:    nil,
			Right:   nil,
			Network: right_network,
		}

		node.Network.Hosts = nil
	} else {
		return fmt.Errorf("Network doesn't support being split.\n")
	}

	return nil
}
