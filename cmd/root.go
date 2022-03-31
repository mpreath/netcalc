package cmd

import (
	"github.com/spf13/cobra"
)

var json_out bool
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "netcalc",
	Short: "Netcalc is a IPv4/IPv6 network calculator",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&json_out, "json", "j", false, "Turns on JSON output for commands")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turns on verbose output for commands")
}

func Execute() error {
	return rootCmd.Execute()
}
