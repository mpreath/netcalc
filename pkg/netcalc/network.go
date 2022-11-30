package netcalc

import (
	"encoding/json"
	"fmt"
	"sort"
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
		Address: ExportAddress(n.Address),
		Mask:    ExportAddress(n.Mask),
		Alias:   (*Alias)(n),
	})
}

func (n *Network) UnmarshalJSON(body []byte) (err error) {
	jsonNetwork := JsonNetwork{}
	if err := json.Unmarshal(body, &jsonNetwork); err != nil {
		return err
	}

	n.Address, err = ParseAddress(jsonNetwork.Address)
	if err != nil {
		return err
	}

	n.Mask, err = ParseAddress(jsonNetwork.Mask)
	if err != nil {
		return err
	}

	return nil
}

func NewNetwork(address uint32, mask uint32) (*Network, error) {
	if !IsValidMask(mask) {
		return nil, fmt.Errorf("network.New: invalid subnet mask")
	}

	return &Network{
		Address: GetNetworkAddress(address, mask),
		Mask:    mask,
	}, nil
}

func (n *Network) Hosts() []*Host {
	var harr []*Host

	for i, address := 0, n.Address+1; address < n.BroadcastAddress(); i, address = i+1, address+1 {
		harr = append(harr, &Host{Address: address, Mask: n.Mask})
	}

	return harr
}

func (n *Network) BroadcastAddress() uint32 {
	return GetBroadcastAddress(n.Address, n.Mask)
}

func (n *Network) HostCount() int {
	return int(n.BroadcastAddress() - n.Address - 1)
}

// SummarizeNetworks takes an array of Network objects and returns a summarized Network.
// Summarization occurs by sorting the array and then finding a common bit boundary that
// a common mask and common network bits can be calculated.
func SummarizeNetworks(networks []*Network) (*Network, error) {

	if len(networks) == 0 {
		return nil, fmt.Errorf("SummarizeNetwork: no networks to summarize")
	} else if len(networks) == 1 {
		return networks[0], nil
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Address < networks[j].Address
	})

	commonBits := networks[0].Address
	var commonMask uint32

	for idx := 1; idx < len(networks); idx++ {
		commonBits = commonBits & networks[idx].Address
		commonMask = GetCommonBitMask(commonBits, networks[idx].Address)
		commonBits = GetNetworkAddress(commonBits, commonMask)
	}

	return &Network{Address: commonBits, Mask: commonMask}, nil

}
