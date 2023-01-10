package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
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
			for _, name := range args {
				if err := db.UpdateTask(name, task.Done); err != nil {
					var nte task.NotFoundErr
					if errors.As(err, &nte) {
						fmt.Printf("ERROR: %s\n", err)
					} else {
						return err
					}
				} else {
					fmt.Printf("Complete task '%s'\n", name)
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(completeCmd)
}
