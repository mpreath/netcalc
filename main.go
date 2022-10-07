package main

import (
	"github.com/mpreath/netcalc/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(2)
	}
}
