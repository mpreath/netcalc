package network

import (
	"fmt"

	"github.com/mpreath/netcalc/ipv4"
	"github.com/mpreath/netcalc/ipv4/host"
)

type Network struct {
	Address          uint32
	Mask             uint32
	BroadcastAddress uint32
	Hosts            []host.Host
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

	network_address := getNetworkAddress(host_address, network_mask)
	broadcast_address := getBroadcastAddress(network_address, network_mask)

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

func getNetworkAddress(address uint32, mask uint32) uint32 {
	return address & mask
}

func getBroadcastAddress(address uint32, mask uint32) uint32 {
	return address | (^mask)
}
