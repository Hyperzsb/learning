package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	completeCmd = &cobra.Command{
		Use:        "complete",
		Aliases:    []string{"c"},
		SuggestFor: []string{"finish"},
		Short:      "Complete a task",
		Long: "Complete a task in 'doing' status, marking it(them) with a 'done' status. " +
			"Multiple tasks names are acceptable",
		GroupID: "task-life-cycle",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("complete")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(completeCmd)
}
