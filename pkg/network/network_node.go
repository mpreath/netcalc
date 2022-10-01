package network

import (
	"fmt"

	"github.com/mpreath/netcalc/pkg/utils"
)

type NetworkNode struct {
	Network  *Network       `json:"network,omitempty"`
	Utilized bool           `json:"-"`
	Subnets  []*NetworkNode `json:"subnets,omitempty"`
}

func (node *NetworkNode) Split() error {
	if len(node.Subnets) > 0 {
		return nil
	}
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

		node.Subnets = append(node.Subnets, &NetworkNode{Network: left_network})

		// right will contain the larger value
		right_network, err := GenerateNetworkFromBits(left_network.BroadcastAddress()+1, new_mask)
		if err != nil {
			return err
		}

		node.Subnets = append(node.Subnets, &NetworkNode{Network: right_network})

	} else {
		return fmt.Errorf("network:Split: network doesn't support being split")
	}

	return nil
}

func GetNetworkCount(node *NetworkNode) int {
	if node == nil {
		return 0
	} else if len(node.Subnets) == 0 {
		return 1
	} else {
		return GetNetworkCount(node.Subnets[0]) + GetNetworkCount(node.Subnets[1])
	}
}

func SplitToVlsmCount(node *NetworkNode, vlsm_count int) error {

	// if network supports requirements
	// and network is not utilized, set utilized to true, return nil
	// if network supports requirements, but is already utilized, then check the
	// next subnet
	// if network doesn't support network
	longest_valid_mask, _ := utils.GetMaskFromBits(30)

	if vlsm_count < 2 {
		return fmt.Errorf("network.SplitToVlsmCount: you must specify at least 2 hosts for count")
	} else {

		// need to determine if this is our recursive base case
		// does this network support our current vlsm_count requirements?
		// 1. check our current host_count
		// 2. look ahead to what the next host count would be
		// 3. if our current host count meets the requirement but our next host count doesn't
		//    then we have found our network
		current_host_count := node.Network.HostCount()
		var lookahead_host_count int
		current_mask_bc := utils.GetBitsInMask(node.Network.Mask)
		lookahead_mask_bc := current_mask_bc + 1
		if lookahead_mask_bc <= 30 {
			// the next split will be a legitimate network
			lookahead_mask, _ := utils.GetMaskFromBits(lookahead_mask_bc)
			lookahead_network, _ := GenerateNetworkFromBits(node.Network.Address, lookahead_mask)
			lookahead_host_count = lookahead_network.HostCount()
		} else {
			lookahead_host_count = 0
		}

		if current_host_count >= vlsm_count && lookahead_host_count < vlsm_count {
			// our current_host_count meets the vlsm count requirements
			// and the next network's count is too small
			// we've found our spot

			// if its not utilized mark it as utilized and return nil
			// if it is utilized return error
			if node.Utilized || len(node.Subnets) > 0 {
				return fmt.Errorf("network.SplitToVlsmCount: network already utilized")
			} else {
				node.Utilized = true
				return nil // no error, base case success
			}
		} else if node.Network.Mask == longest_valid_mask {
			return fmt.Errorf("network.SplitToVlsmCount: network can't support that many subnetworks")
		}

		err := node.Split()
		if err != nil {
			return err
		}

		err = SplitToVlsmCount(node.Subnets[0], vlsm_count)
		if err != nil {
			err = SplitToVlsmCount(node.Subnets[1], vlsm_count)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func NetworkNodeToArray(node *NetworkNode) []*Network {
	var narr []*Network
	if len(node.Subnets) > 0 {
		narr = append(narr, NetworkNodeToArray(node.Subnets[0])...)
		narr = append(narr, NetworkNodeToArray(node.Subnets[1])...)
	} else if len(node.Subnets) == 0 {
		narr = append(narr, node.Network)
	}
	return narr
}
