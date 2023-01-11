package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
)

var (
	undoComplete bool
	completeCmd  = &cobra.Command{
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
				var status task.Status
				if undoComplete {
					status = task.Doing
				} else {
					status = task.Done
				}

				if err := db.UpdateTask(name, status); err != nil {
					var nte task.NotFoundErr
					if errors.As(err, &nte) {
						fmt.Printf("ERROR: %s\n", err)
					} else {
						return err
					}
				}

				if undoComplete {
					fmt.Printf("Doing task '%s'\n", name)
				} else {
					fmt.Printf("Done task '%s'\n", name)
				}
			}
			return nil
		},
	}
)

func init() {
	completeCmd.Flags().BoolVarP(&undoComplete, "undo", "u", false, "change status of a 'done' task back to 'doing'")

	rootCmd.AddCommand(completeCmd)
}
