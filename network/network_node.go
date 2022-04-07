package network

import (
	"fmt"
	"math"

	"github.com/mpreath/netcalc/utils"
)

type NetworkNode struct {
	Network *Network       `json:"network,omitempty"`
	Parent  *NetworkNode   `json:"-"`
	Subnets []*NetworkNode `json:"subnets,omitempty"`
}

func (node *NetworkNode) Split() error {
	bc := utils.GetBitsInMask(node.Network.Mask) + 1
	if bc < 31 {
		new_mask, err := utils.GetMaskFromBits(bc)
		if err != nil {
			return err
		}
		// left will contain the lower value
		left_network, err := GenerateNetworkFromBits(node.Network.Address, new_mask)
		if err != nil {
			return err
		}

		node.Subnets = append(node.Subnets, &NetworkNode{Parent: node, Network: left_network})

		// right will contain the larger value
		right_network, err := GenerateNetworkFromBits(left_network.BroadcastAddress+1, new_mask)
		if err != nil {
			return err
		}

		node.Subnets = append(node.Subnets, &NetworkNode{Parent: node, Network: right_network})

		// no usable hosts in this network
		node.Network.Hosts = nil
	} else {
		return fmt.Errorf("network doesn't support being split")
	}

	return nil
}

func SplitToHostCount(node *NetworkNode, host_count int) error {
	current_mask_bc := utils.GetBitsInMask(node.Network.Mask)
	if current_mask_bc >= 30 {
		// this is the longest mask we support
		return nil
	}
	current_bc := 32 - current_mask_bc
	current_hc := int(math.Pow(2, float64(current_bc)))
	future_bc := current_bc - 1 // need to look ahead into the future
	future_hc := int(math.Pow(2, float64(future_bc)))

	if current_hc >= host_count && future_hc < host_count {
		// this is our recursive base case
		return nil
	} else if current_hc < host_count {
		// requirements too large, raise an error
		return fmt.Errorf("network can't support that many hosts")
	} else {
		err := node.Split()
		if err != nil {
			return err
		}
		err = SplitToHostCount(node.Subnets[0], host_count)
		if err != nil {
			return err
		}
		err = SplitToHostCount(node.Subnets[1], host_count)
		if err != nil {
			return err
		}
	}

	return nil
}

func SplitToNetCount(node *NetworkNode, net_count int) error {
	longest_valid_mask, _ := utils.GetMaskFromBits(30)
	if net_count <= 0 {
		// this is our recursive base case
		return nil
	} else if node.Network.Mask == longest_valid_mask {
		// can't split any more
		return fmt.Errorf("network can't support that many subnetworks")
	} else {
		err := node.Split()
		if err != nil {
			return err
		}
		err = SplitToNetCount(node.Subnets[0], net_count-2)
		if err != nil {
			return err
		}
		err = SplitToNetCount(node.Subnets[1], net_count-2)
		if err != nil {
			return err
		}
	}

	return nil
}
