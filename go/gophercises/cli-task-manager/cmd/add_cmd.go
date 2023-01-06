package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:        "add",
		Aliases:    []string{"a"},
		SuggestFor: []string{"new", "create"},
		Short:      "Add a new task",
		Long: "Add and start to trace a new task, marking the task with a 'todo' status. " +
			"Multiple tasks names are acceptable",
		GroupID: "task-life-cycle",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("add")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
