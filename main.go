package main

import (
	"fmt"
	"netcalc/ipv4/host"
)

func main() {
	host := host.Host{
		Address: 15,
	}

	fmt.Println(host.Address)
}
