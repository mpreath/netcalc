package host

import (
	"encoding/json"

	"github.com/mpreath/netcalc/utils"
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
		Address: utils.Itodd(h.Address),
		Mask:    utils.Itodd(h.Mask),
	})
}

func GenerateHost(address string, mask string) (*Host, error) {
	host_address, err := utils.Ddtoi(address)
	if err != nil {
		return nil, err
	}

	host_mask, err := utils.Ddtoi(mask)
	if err != nil {
		return nil, err
	}

	return &Host{Address: host_address, Mask: host_mask}, nil
}
