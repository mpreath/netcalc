// Package network provides types and methods for working with
// TCP/IP networks. Each network has an Address (uint32) and Mask (uint32)
// that represent a give IPv4 network.
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

type JsonNetwork struct {
	Address string `json:"address"`
	Mask    string `json:"mask"`
}

func (n *Network) MarshalJSON() ([]byte, error) {
	type Alias Network
	return json.Marshal(&struct {
		Address string `json:"address"`
		Mask    string `json:"mask"`
		*Alias
	}{
		Address: utils.ExportAddress(n.Address),
		Mask:    utils.ExportAddress(n.Mask),
		Alias:   (*Alias)(n),
	})
}

func (n *Network) UnmarshalJSON(body []byte) (err error) {
	jsonNetwork := JsonNetwork{}
	if err := json.Unmarshal(body, &jsonNetwork); err != nil {
		return err
	}

	n.Address, err = utils.ParseAddress(jsonNetwork.Address)
	if err != nil {
		return err
	}

	n.Mask, err = utils.ParseAddress(jsonNetwork.Mask)
	if err != nil {
		return err
	}

	return nil
}

func New(address uint32, mask uint32) (*Network, error) {
	if !utils.IsValidMask(mask) {
		return nil, fmt.Errorf("network.New: invalid subnet mask")
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
