package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	force bool

	removeCmd = &cobra.Command{
		Use:        "remove",
		Aliases:    []string{"rm"},
		SuggestFor: []string{"delete", "prune"},
		Short:      "Remove a task",
		Long: "Remove a task, no matter which status it has. " +
			"Multiple tasks names are acceptable",
		GroupID: "task-life-cycle",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("remove")
			return nil
		},
	}
)

func init() {
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "remove the task without asking for consent")

	rootCmd.AddCommand(removeCmd)
}
