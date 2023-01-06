package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	FirstName string
	LastName  string

	flagCmd = &cobra.Command{
		Use:   "flag",
		Short: "use flags to change behaviors",
		Long:  "use flags to change behaviors of this app to see how Cobra parse and handle them",
		Run: func(cmd *cobra.Command, args []string) {
			if Verbose {
				fmt.Printf("Hello %s %s! Nice to meet you!\n", FirstName, LastName)
			} else {
				fmt.Printf("Hello %s %s\n", FirstName, LastName)
			}
		},
	}
)

func init() {
	flagCmd.Flags().StringVarP(&FirstName, "first", "f", "Harry", "your first name")
	flagCmd.Flags().StringVarP(&LastName, "last", "l", "Potter", "your last name")
	_ = flagCmd.MarkFlagRequired("first")
	rootCmd.AddCommand(flagCmd)
}
