package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gophercises/taskmanager/db"
	"gophercises/taskmanager/task"
)

var (
	todo  bool
	doing bool
	done  bool
	all   bool

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List tasks",
		Long:    "List todo tasks by default, or you can specify tasks with different status",
		GroupID: "general",
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, err := db.ListTask()
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				fmt.Println("No task available")
			}

			for _, t := range tasks {
				switch t.Status {
				case task.Todo:
					if (!todo && !doing && !done) || todo || all {
						fmt.Println(t.ToString())
					}
				case task.Doing:
					if doing || all {
						fmt.Println(t.ToString())
					}
				case task.Done:
					if done || all {
						fmt.Println(t.ToString())
					}
				}
			}

			return nil
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&todo, "todo", "", false, "list todo tasks")
	listCmd.Flags().BoolVarP(&doing, "doing", "", false, "list doing tasks")
	listCmd.Flags().BoolVarP(&done, "done", "", false, "list done tasks")
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "list all tasks")

	rootCmd.AddCommand(listCmd)
}
