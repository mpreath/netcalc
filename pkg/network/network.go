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

func (n *Network) Hosts() []*host.Host {
	var harr []*host.Host

	for i, address := 0, n.Address+1; address < n.BroadcastAddress(); i, address = i+1, address+1 {
		harr = append(harr, &host.Host{Address: address, Mask: n.Mask})
	}

	return harr
}

func (n *Network) BroadcastAddress() uint32 {
	return utils.GetBroadcastAddress(n.Address, n.Mask)
}

func (n *Network) HostCount() int {
	return int(n.BroadcastAddress() - n.Address - 1)
}
