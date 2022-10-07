package main

import (
	"github.com/spf13/cobra"
)

var JSON_FLAG bool
var VERBOSE_FLAG bool

var rootCmd = &cobra.Command{
	Use:   "netcalc",
	Short: "Netcalc is a IPv4/IPv6 network calculator",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&JSON_FLAG, "json", "j", false, "Turns on JSON output for commands")
	rootCmd.PersistentFlags().BoolVarP(&VERBOSE_FLAG, "verbose", "v", false, "Turns on verbose output for commands")
}

func Execute() error {
	return rootCmd.Execute()
}
