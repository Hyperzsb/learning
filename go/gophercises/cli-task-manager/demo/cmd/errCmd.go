package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var errCmd = &cobra.Command{
	Use:   "err",
	Short: "cause an error intentionally",
	Long:  "cause an error intentionally in this app to see how Cobra handle it",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("this is an intentionally caused error")
	},
}

func init() {
	rootCmd.AddCommand(errCmd)
}
