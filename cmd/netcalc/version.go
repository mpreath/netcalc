package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VERSION = "0.2"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of netcalc",
	Run: func(cmd *cobra.Command, args []string) {
		printVersionInformation()
	},
}

func printVersionInformation() {
	fmt.Printf("netcalc %s\n", VERSION)
}
