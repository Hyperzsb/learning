package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
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
			for _, name := range args {
				if !force {
					fmt.Printf("Remove task '%s'? (y/n, default n) ", name)
					confirm := ""
					_, _ = fmt.Scanf("%s", &confirm)
					if confirm != "y" {
						continue
					}
				}

				if err := db.RemoveTask(name); err != nil {
					var nte task.NotFoundErr
					if errors.As(err, &nte) {
						fmt.Printf("ERROR: %s\n", err)
					} else {
						return err
					}
				} else {
					fmt.Printf("Task '%s' removed\n", name)
				}
			}
			return nil
		},
	}
)

func init() {
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "remove the task without asking for consent")

	rootCmd.AddCommand(removeCmd)
}
