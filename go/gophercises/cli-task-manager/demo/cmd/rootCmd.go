package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Verbose bool

	rootCmd = &cobra.Command{
		Use:   "cobra-demo",
		Short: "A demo for Cobra",
		Long:  "This is a tiny demo for learning how to use Cobra to build a CLI app.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
		},
	}
)

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "turn on verbose mode")
	return rootCmd.Execute()
}
