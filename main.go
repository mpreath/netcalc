package main

import (
	"fmt"
	"netcalc/ipv4"
	"netcalc/ipv4/network"
)

func main() {
	n, err := network.GenerateNetwork("192.168.1.90", "255.255.255.252")
	if err != nil {
		fmt.Print(err)
		return
	}
	n_dd_address, _ := ipv4.Itodd(n.Address)
	n_dd_mask, _ := ipv4.Itodd(n.Mask)
	n_dd_bcast, _ := ipv4.Itodd(n.BroadcastAddress)
	fmt.Printf("Network: { Address: \"%s\", Mask: \"%s\", BroadcastAddress: \"%s\", Host Count: %d }\n", n_dd_address, n_dd_mask, n_dd_bcast, len(n.Hosts))
}
