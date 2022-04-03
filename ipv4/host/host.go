package host

import (
	"encoding/json"

	"github.com/mpreath/netcalc/ipv4"
)

type Host struct {
	Address uint32 `json:"address"`
	Mask    uint32 `json:"mask"`
}

func (h *Host) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Address string `json:"address"`
		Mask    string `json:"mask"`
	}{
		Address: ipv4.Itodd(h.Address),
		Mask:    ipv4.Itodd(h.Mask),
	})
}

func GenerateHost(address string, mask string) (*Host, error) {
	host_address, err := ipv4.Ddtoi(address)
	if err != nil {
		return nil, err
	}

	host_mask, err := ipv4.Ddtoi(mask)
	if err != nil {
		return nil, err
	}

	return &Host{Address: host_address, Mask: host_mask}, nil
}
