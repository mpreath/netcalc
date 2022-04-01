package network

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/host"
)

type Network struct {
	Address          uint32      `json:"address"`
	Mask             uint32      `json:"mask"`
	BroadcastAddress uint32      `json:"broadcast"`
	Hosts            []host.Host `json:"hosts,omitempty"`
}

func (n *Network) MarshalJSON() ([]byte, error) {
	type Alias Network
	return json.Marshal(&struct {
		Address          string `json:"address"`
		Mask             string `json:"mask"`
		BroadcastAddress string `json:"broadcast"`
		*Alias
	}{
		Address:          ipv4.Itodd(n.Address),
		Mask:             ipv4.Itodd(n.Mask),
		BroadcastAddress: ipv4.Itodd(n.BroadcastAddress),
		Alias:            (*Alias)(n),
	})
}

func GenerateNetwork(address string, mask string) (*Network, error) {
	host_address, err := ipv4.Ddtoi(address)
	if err != nil {
		return nil, err
	}

	network_mask, err := ipv4.Ddtoi(mask)
	if err != nil {
		return nil, err
	} else if !ipv4.IsValidMask(network_mask) {
		return nil, fmt.Errorf("ipv4.network.GenerateNetwork: invalid subnet mask")
	}

	network_address := ipv4.GetNetworkAddress(host_address, network_mask)
	broadcast_address := ipv4.GetBroadcastAddress(network_address, network_mask)

	count := broadcast_address - network_address - 1

	network := Network{
		Address:          network_address,
		Mask:             network_mask,
		BroadcastAddress: broadcast_address,
		Hosts:            make([]host.Host, count),
	}

	for i, n := 0, network_address+1; n < broadcast_address; i, n = i+1, n+1 {
		network.Hosts[i].Address = n
		network.Hosts[i].Mask = network_mask
	}

	return &network, nil
}

func GenerateNetworkFromBits(address uint32, mask uint32) (*Network, error) {

	network_address := ipv4.GetNetworkAddress(address, mask)
	broadcast_address := ipv4.GetBroadcastAddress(address, mask)

	count := broadcast_address - network_address - 1

	network := Network{
		Address:          network_address,
		Mask:             mask,
		BroadcastAddress: broadcast_address,
		Hosts:            make([]host.Host, count),
	}

	for i, n := 0, network_address+1; n < broadcast_address; i, n = i+1, n+1 {
		network.Hosts[i].Address = n
		network.Hosts[i].Mask = mask
	}
	return &network, nil
}
