package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
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
			for _, name := range args {
				if err := db.UpdateTask(name, task.Doing); err != nil {
					var nte task.NotFoundErr
					if errors.As(err, &nte) {
						fmt.Printf("ERROR: %s\n", err)
					} else {
						return err
					}
				} else {
					fmt.Printf("Do task '%s'\n", name)
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(doCmd)
}
