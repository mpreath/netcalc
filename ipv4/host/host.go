package host

import "netcalc/ipv4"

type Host struct {
	Address uint32
	Mask    uint32
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
