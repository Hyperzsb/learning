package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	doCmd = &cobra.Command{
		Use:        "do",
		Aliases:    []string{"d"},
		SuggestFor: []string{"start", "begin"},
		Short:      "Start to do a task",
		Long: "Start to do a task in 'todo' status, marking it with a 'doing' status. " +
			"Multiple tasks names are acceptable",
		GroupID: "task-life-cycle",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("do")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(doCmd)
}
