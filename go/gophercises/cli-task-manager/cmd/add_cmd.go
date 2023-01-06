package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
	"os"
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
			for _, name := range args {
				// TODO: validate user input
				var newTask task.Task
				newTask.Name = name
				newTask.Status = task.Todo

				fmt.Printf("Please prvoide description of task '%s': \n> ", newTask.Name)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				newTask.Desc = scanner.Text()

				err := db.AddTask(newTask)
				if err != nil {
					// TODO: figure the best way to determine the actual type of an error
					//var dev task.DuplicateErr
					//dev.Name = newTask.Name
					//if errors.Is(err, dev) {
					//	fmt.Println("!", err)
					//} else {
					//	fmt.Println("?", err)
					//}

					var dep *task.DuplicateErr
					if errors.As(err, &dep) {
						fmt.Println(err)
					} else {
						return err
					}

					//if e, ok := err.(*task.DuplicateErr); ok {
					//	fmt.Println("!", e)
					//} else {
					//	fmt.Println("?", e)
					//}
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
