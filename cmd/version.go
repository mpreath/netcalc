package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	fmt.Println("netcalc 0.1")
}
