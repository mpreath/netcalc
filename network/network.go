package network

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/utils"
)

type Network struct {
	Address          uint32 `json:"address"`
	Mask             uint32 `json:"mask"`
	BroadcastAddress uint32 `json:"broadcast"`
	HostCount        uint   `json:"host_count"`
}

func (n *Network) MarshalJSON() ([]byte, error) {
	type Alias Network
	return json.Marshal(&struct {
		Address          string `json:"address"`
		Mask             string `json:"mask"`
		BroadcastAddress string `json:"broadcast"`
		*Alias
	}{
		Address:          utils.Itodd(n.Address),
		Mask:             utils.Itodd(n.Mask),
		BroadcastAddress: utils.Itodd(n.BroadcastAddress),
		Alias:            (*Alias)(n),
	})
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

	network_address := utils.GetNetworkAddress(host_address, network_mask)
	broadcast_address := utils.GetBroadcastAddress(network_address, network_mask)

	// count := broadcast_address - network_address - 1

	network := Network{
		Address:          network_address,
		Mask:             network_mask,
		BroadcastAddress: broadcast_address,
		HostCount:        uint(broadcast_address - network_address - 1),
	}

	// for i, n := 0, network_address+1; n < broadcast_address; i, n = i+1, n+1 {
	// 	network.Hosts[i].Address = n
	// 	network.Hosts[i].Mask = network_mask
	// }

	return &network, nil
}

func GenerateNetworkFromBits(address uint32, mask uint32) (*Network, error) {

	network_address := utils.GetNetworkAddress(address, mask)
	broadcast_address := utils.GetBroadcastAddress(address, mask)

	network := Network{
		Address:          network_address,
		Mask:             mask,
		BroadcastAddress: broadcast_address,
		HostCount:        uint(broadcast_address - network_address - 1),
	}

	return &network, nil
}
