package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var argCmd = &cobra.Command{
	Use:   "arg",
	Short: "use some arguments to change behaviors",
	Long:  "use some arguments to change behaviors of this app to see how Cobra parse and handle them",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			fmt.Printf("Hello, Mr./Miss. %s\n", args[0])
		} else {
			fmt.Printf("Hello, Mr./Miss. %s %s\n", args[0], args[1])
		}
	},
}

func init() {
	rootCmd.AddCommand(argCmd)
}
