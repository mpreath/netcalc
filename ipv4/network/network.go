package network

import "netcalc/ipv4/host"

type Network struct {
	Address uint32
	Mask    uint32
	Hosts   []host.Host
}
