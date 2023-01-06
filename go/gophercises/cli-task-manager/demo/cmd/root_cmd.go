package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Verbose bool

	rootCmd = &cobra.Command{
		Use:   "cobra-demo",
		Short: "A demo for Cobra",
		Long:  "This is a tiny demo for learning how to use Cobra to build a CLI app.",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "turn on verbose mode")
}

func Execute() error {
	return rootCmd.Execute()
}
