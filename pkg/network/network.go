package network

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/pkg/host"
	"github.com/mpreath/netcalc/pkg/utils"
)

type Network struct {
	Address uint32 `json:"address"`
	Mask    uint32 `json:"mask"`
}

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

func New(address uint32, mask uint32) (*Network, error) {
	if !utils.IsValidMask(mask) {
		return nil, fmt.Errorf("network.GenerateNetwork: invalid subnet mask")
	}

	return &Network{
		Address: utils.GetNetworkAddress(address, mask),
		Mask:    mask,
	}, nil
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

func (n *Network) BroadcastAddress() uint32 {
	return utils.GetBroadcastAddress(n.Address, n.Mask)
}

func (n *Network) HostCount() int {
	return int(n.BroadcastAddress() - n.Address - 1)
}
