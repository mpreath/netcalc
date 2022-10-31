package host

import (
	"encoding/json"
	"fmt"

	"github.com/mpreath/netcalc/pkg/utils"
)

// TODO: add doc related comments
type Host struct {
	Address uint32 `json:"address"`
	Mask    uint32 `json:"mask"`
}

func (h *Host) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Address string `json:"address"`
		Mask    string `json:"mask"`
	}{
		Address: utils.ExportAddress(h.Address),
		Mask:    utils.ExportAddress(h.Mask),
	})
}

func New(address uint32, mask uint32) (*Host, error) {

	if !utils.IsValidMask(mask) {
		return nil, fmt.Errorf("host.New: invalid subnet mask")
	}

	return &Host{Address: address, Mask: mask}, nil
}
