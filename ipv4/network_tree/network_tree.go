package network_tree

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/network"
)

type NetworkNode struct {
	Parent *NetworkNode
	Left   *NetworkNode
	Right  *NetworkNode
	// wrap a network in each node struct
	*network.Network
}

func (node *NetworkNode) Print() {
	s, _ := json.MarshalIndent(node, " ", "    ")
	fmt.Printf("%s\n", s)
}

func (node *NetworkNode) Split() error {
	bc := ipv4.GetBitsInMask(node.Mask) + 1
	if bc < 31 {
		new_mask, err := ipv4.GetMaskFromBits(bc)
		if err != nil {
			return err
		}
		// left will contain the lower value
		left_network, err := network.GenerateNetworkFromBits(node.Address, new_mask)
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

		node.Hosts = nil
	} else {
		return fmt.Errorf("Network doesn't support being split.\n")
	}

	return nil
}
