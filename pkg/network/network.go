package network

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/mpreath/netcalc/pkg/host"
	"github.com/mpreath/netcalc/pkg/utils"
)

type Network struct {
	Address uint32 `json:"address"` // 32 bits
	Mask    uint32 `json:"mask"`    // 32 bits
} // 8 bytes
// 2^30 = 1073741824
// 8 bytes * 1073741824 = 8589934592 bytes
// 8589934592 bytes * 1 KB/1000B * 1MB/1000KB * 1GB/1000MB = 8.59GB

func (n *Network) MarshalJSON() ([]byte, error) {
	type Alias Network
	return json.Marshal(&struct {
		Address string `json:"address"`
		Mask    string `json:"mask"`
		*Alias
	}{
		Address: utils.Itodd(n.Address),
		Mask:    utils.Itodd(n.Mask),
		Alias:   (*Alias)(n),
	})
}

func (n *Network) BroadcastAddress() uint32 {
	return utils.GetBroadcastAddress(n.Address, n.Mask)
}

func (n *Network) HostCount() int {
	return int(n.BroadcastAddress() - n.Address - 1)
}

func GenerateNetwork(address string, mask string) (*Network, error) {
	host_address, err := utils.Ddtoi(address)
	if err != nil {
		return nil, err
	}

	network_mask, err := utils.Ddtoi(mask)
	if err != nil {
		return nil, err
	} else if !utils.IsValidMask(network_mask) {
		return nil, fmt.Errorf("network.GenerateNetwork: invalid subnet mask")
	}

	network := Network{
		Address: utils.GetNetworkAddress(host_address, network_mask),
		Mask:    network_mask,
	}

	return &network, nil
}

func GenerateNetworkFromBits(address uint32, mask uint32) (*Network, error) {

	network := Network{
		Address: utils.GetNetworkAddress(address, mask),
		Mask:    mask,
	}

	return &network, nil
}

func GetHosts(network *Network) []*host.Host {
	var harr []*host.Host

	for i, n := 0, network.Address+1; n < network.BroadcastAddress(); i, n = i+1, n+1 {
		harr = append(harr, &host.Host{Address: n, Mask: network.Mask})
	}

	return harr
}

func (n *Network) Split() (*Network, *Network, error) {
	var left_network, right_network *Network
	bc := utils.GetBitsInMask(n.Mask) + 1
	if bc < 31 {
		new_mask, err := utils.GetMaskFromBits(bc)
		if err != nil {
			return nil, nil, err
		}
		// left will contain the lower value
		left_network, err = GenerateNetworkFromBits(n.Address, new_mask)
		if err != nil {
			return nil, nil, err
		}
		// right will contain the larger value
		right_network, err = GenerateNetworkFromBits(left_network.BroadcastAddress()+1, new_mask)
		if err != nil {
			return nil, nil, err
		}

	} else {
		return nil, nil, fmt.Errorf("network:Split: network doesn't support being split")
	}

	return left_network, right_network, nil
}

func SplitToHostCount(net *Network, host_count int) ([]*Network, error) {
	current_mask_bc := utils.GetBitsInMask(net.Mask)
	if current_mask_bc >= 30 {
		// this is the longest mask we support
		return []*Network{net}, nil
	}
	current_bc := 32 - current_mask_bc
	current_hc := int(math.Pow(2, float64(current_bc)))
	future_bc := current_bc - 1 // need to look ahead into the future
	future_hc := int(math.Pow(2, float64(future_bc)))

	if current_hc >= host_count && future_hc < host_count {
		// this is our recursive base case
		return []*Network{net}, nil
	} else if current_hc < host_count {
		// requirements too large, raise an error
		return nil, fmt.Errorf("network.SplitToHostCount: network can't support that many hosts")
	} else {
		net1, net2, err := net.Split()
		if err != nil {
			return nil, err
		}
		res1, err := SplitToHostCount(net1, host_count)
		if err != nil {
			return nil, err
		}
		res2, err := SplitToHostCount(net2, host_count)
		if err != nil {
			return nil, err
		}
		// FIXME: Not sure appending is the most efficient way to do this, BUT might be
		return append(res1, res2...), nil
	}
}

func SplitToNetCount(net *Network, net_count int) ([]*Network, error) {
	longest_valid_mask, _ := utils.GetMaskFromBits(30)
	if net_count <= 0 {
		// this is our recursive base case
		return []*Network{net}, nil
	} else if net.Mask == longest_valid_mask {
		// can't split any more
		return nil, fmt.Errorf("network.SplitToNetCount: network can't support that many subnetworks")
	} else {
		net1, net2, err := net.Split()
		if err != nil {
			return nil, err
		}
		res1, err := SplitToNetCount(net1, net_count-2)
		if err != nil {
			return nil, err
		}
		res2, err := SplitToNetCount(net2, net_count-2)
		if err != nil {
			return nil, err
		}
		return append(res1, res2...), nil

	}
}
