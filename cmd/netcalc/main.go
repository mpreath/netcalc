package main

import (
	"github.com/mpreath/netcalc/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
