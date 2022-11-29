// Package host provides a Host type consisting of two 32-bit unsigned integers
// representing a 32-bit IP Address (Address) and a 32-bit subnet mask (Mask).
package netcalc

import (
	"encoding/json"
	"fmt"
)

// Host type consisting of two 32-bit unsigned integers
// representing a 32-bit IP Address (Address) and a 32-bit subnet mask (Mask).
type Host struct {
	Address uint32 `json:"address"`
	Mask    uint32 `json:"mask"`
}

func (h *Host) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Address string `json:"address"`
		Mask    string `json:"mask"`
	}{
		Address: ExportAddress(h.Address),
		Mask:    ExportAddress(h.Mask),
	})
}

// NewHost initializes and returns a Host based on the address and mask arguments.
// It returns an error if the mask is invalid.
func NewHost(address uint32, mask uint32) (*Host, error) {

	if !IsValidMask(mask) {
		return nil, fmt.Errorf("host.New: invalid subnet mask")
	}

	return &Host{Address: address, Mask: mask}, nil
}
