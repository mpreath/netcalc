// Package networknode provides types and methods for working with
// TCP/IP networks for creating a subnetworks and summarizing subnetworks.
package networknode

import (
	"fmt"
	"github.com/mpreath/netcalc/pkg/network"
	"math"

	"github.com/mpreath/netcalc/pkg/utils"
)

// NetworkNode is a wrapper around the network.Network struct adding
// a Utilized field (used internally) and an array of NetworkNode pointers
// to build the subnet binary tree.
type NetworkNode struct {
	Network  *network.Network `json:"network,omitempty"`
	Utilized bool             `json:"-"`
	Subnets  []*NetworkNode   `json:"subnets,omitempty"`
}

func New(n *network.Network) *NetworkNode {
	return &NetworkNode{
		Network: n,
	}
}

// Split adds two subnetworks (subnets) to the Subnets array.
// The subnets created have a mask that is extended by 1-bit
// effectively splitting the original network.
func (node *NetworkNode) Split() error {
	if len(node.Subnets) > 0 {
		return nil
	}
	bc := utils.GetBitsInMask(node.Network.Mask) + 1
	if bc < 31 {
		newMask, err := utils.GetMaskFromBits(bc)
		if err != nil {
			return err
		}
		// left will contain the lower value
		leftNetwork, err := network.New(node.Network.Address, newMask)
		if err != nil {
			return err
		}

		node.Subnets = append(node.Subnets, &NetworkNode{Network: leftNetwork})

		// right will contain the larger value
		rightNetwork, err := network.New(leftNetwork.BroadcastAddress()+1, newMask)
		if err != nil {
			return err
		}

		node.Subnets = append(node.Subnets, &NetworkNode{Network: rightNetwork})

	} else {
		return fmt.Errorf("network:Split: network doesn't support being split")
	}

	return nil
}

func (node *NetworkNode) NetworkCount() int {
	if node == nil {
		return 0
	} else if len(node.Subnets) == 0 {
		return 1
	} else {
		return node.Subnets[0].NetworkCount() + node.Subnets[1].NetworkCount()
	}
}

// SplitToHostCount provided with a NetworkNode and hostCount, recursively calls Split on the
// node(s) until the host count of the networks created support the hostCount provided.
func SplitToHostCount(node *NetworkNode, hostCount int) error {

	valid, err := ValidForHostCount(node.Network, hostCount)
	if err != nil {
		return err
	}
	if valid {
		return nil // success
	} else {
		err := node.Split()
		if err != nil {
			return err
		}
		err = SplitToHostCount(node.Subnets[0], hostCount)
		if err != nil {
			return err
		}
		err = SplitToHostCount(node.Subnets[1], hostCount)
		if err != nil {
			return err
		}

		return nil
	}
}

// ValidForHostCount provided a Network and hostCount returns true if the current network's
// host count supports the provided hostCount and a further split network would not.
func ValidForHostCount(n *network.Network, hostCount int) (bool, error) {

	currentMaskBc := utils.GetBitsInMask(n.Mask)
	currentBc := 32 - currentMaskBc
	currentHc := int(math.Pow(2, float64(currentBc)))
	futureBc := currentBc - 1 // need to look ahead into the future
	futureHc := int(math.Pow(2, float64(futureBc)))

	if currentHc >= hostCount && futureHc < hostCount {
		// this is our recursive base case
		return true, nil
	} else if currentHc < hostCount {
		// requirements too large, raise an error
		return false, fmt.Errorf("network.SplitToHostCount: network can't support that many hosts")
	} else if currentMaskBc >= 30 {
		return true, nil
	}

	return false, nil
}

// SplitToNetCount provided with a NetworkNode and netCount, recursively calls Split on the
// node(s) until the number of networks created is equal to or greater than the netCount argument.
func SplitToNetCount(node *NetworkNode, netCount int) error {
	longestValidMask, _ := utils.GetMaskFromBits(30)
	if netCount <= 0 {
		// this is our recursive base case
		return nil
	} else if node.Network.Mask == longestValidMask {
		// can't split any more
		return fmt.Errorf("network.SplitToNetCount: network can't support that many subnetworks")
	} else {
		err := node.Split()
		if err != nil {
			return err
		}
		err = SplitToNetCount(node.Subnets[0], netCount-2)
		if err != nil {
			return err
		}
		err = SplitToNetCount(node.Subnets[1], netCount-2)
		if err != nil {
			return err
		}
	}

	return nil
}

// SplitToVlsmCount provided with a NetworkNode and vlsmCount value (hostCount), recursively calls
// Split on the node until the host count of the split network is valid for the vlsmCount.
// This differs from SplitToHostCount because it recursively goes down one branch of the tree rather than
// all branches creating an unbalanced binary tree.
func SplitToVlsmCount(node *NetworkNode, vlsmCount int) error {

	// if network supports requirements
	// and network is not utilized, set utilized to true, return nil
	// if network supports requirements, but is already utilized, then check the
	// next subnet
	// if network doesn't support network
	longestValidMask, _ := utils.GetMaskFromBits(30)

	if vlsmCount < 2 {
		return fmt.Errorf("network.SplitToVlsmCount: you must specify at least 2 hosts for count")
	} else {

		// need to determine if this is our recursive base case
		// does this network support our current vlsmCount requirements?
		// 1. check our current hostCount
		// 2. look ahead to what the next host count would be
		// 3. if our current host count meets the requirement but our next host count doesn't
		//    then we have found our network
		currentHostCount := node.Network.HostCount()
		var lookaheadHostCount int
		currentMaskBc := utils.GetBitsInMask(node.Network.Mask)
		lookaheadMaskBc := currentMaskBc + 1
		if lookaheadMaskBc <= 30 {
			// the next split will be a legitimate network
			lookaheadMask, _ := utils.GetMaskFromBits(lookaheadMaskBc)
			lookaheadNetwork, _ := network.New(node.Network.Address, lookaheadMask)
			lookaheadHostCount = lookaheadNetwork.HostCount()
		} else {
			lookaheadHostCount = 0
		}

		if currentHostCount >= vlsmCount && lookaheadHostCount < vlsmCount {
			// our currentHostCount meets the vlsm count requirements
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
		} else if node.Network.Mask == longestValidMask {
			return fmt.Errorf("network.SplitToVlsmCount: network can't support that many subnetworks")
		}

		err := node.Split()
		if err != nil {
			return err
		}

		err = SplitToVlsmCount(node.Subnets[0], vlsmCount)
		if err != nil {
			err = SplitToVlsmCount(node.Subnets[1], vlsmCount)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// Flatten returns the leaf nodes of the binary tree as an array of Networks.
func (node *NetworkNode) Flatten() []*network.Network {
	var networkList []*network.Network
	if len(node.Subnets) == 0 {
		return append(networkList, node.Network)
	} else {
		networkList = append(networkList, node.Subnets[0].Flatten()...)
		networkList = append(networkList, node.Subnets[1].Flatten()...)
	}
	return networkList
}

// FlattenUtilized returns the leaf nodes with Utilized set to true of the binary tree as an array of Networks.
func (node *NetworkNode) FlattenUtilized() []*network.Network {
	var networkList []*network.Network
	if node.Utilized {
		return append(networkList, node.Network)
	} else if len(node.Subnets) > 0 {
		networkList = append(networkList, node.Subnets[0].FlattenUtilized()...)
		networkList = append(networkList, node.Subnets[1].FlattenUtilized()...)
	}
	return networkList
}
