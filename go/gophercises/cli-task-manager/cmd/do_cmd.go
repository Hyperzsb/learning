package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
)

var (
	undoDo bool
	doCmd  = &cobra.Command{
		Use:        "do",
		Aliases:    []string{"d"},
		SuggestFor: []string{"start", "begin"},
		Short:      "Start to do a task",
		Long: "Start to do a task in 'todo' status, marking it with a 'doing' status. " +
			"Multiple tasks names are acceptable",
		GroupID: "task-life-cycle",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				var status task.Status
				if undoDo {
					status = task.Todo
				} else {
					status = task.Doing
				}

				if err := db.UpdateTask(name, status); err != nil {
					var nte task.NotFoundErr
					if errors.As(err, &nte) {
						fmt.Printf("ERROR: %s\n", err)
					} else {
						return err
					}
				}

				if undoDo {
					fmt.Printf("Todo task '%s'\n", name)
				} else {
					fmt.Printf("Doing task '%s'\n", name)
				}
			}
			return nil
		},
	}
)

func init() {
	doCmd.Flags().BoolVarP(&undoDo, "undo", "u", false, "change status of a 'doing' task back to 'todo'")

	rootCmd.AddCommand(doCmd)
}
