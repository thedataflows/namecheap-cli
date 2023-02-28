/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
