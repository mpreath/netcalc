package main

import (
	"github.com/spf13/cobra"
)

var JsonFlag bool
var VerboseFlag bool

var rootCmd = &cobra.Command{
	Use:   "netcalc",
	Short: "Netcalc is a IPv4/IPv6 network calculator",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&JsonFlag, "json", "j", false, "Turns on JSON output for commands")
	rootCmd.PersistentFlags().BoolVarP(&VerboseFlag, "verbose", "v", false, "Turns on verbose output for commands")
}

func Execute() error {
	return rootCmd.Execute()
}
