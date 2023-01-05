package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version info",
	Long:  "print the version info of this app and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
