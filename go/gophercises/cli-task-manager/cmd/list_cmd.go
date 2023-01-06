package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
			fmt.Println("list")
			return nil
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&todo, "todo", "", true, "list todo tasks")
	listCmd.Flags().BoolVarP(&doing, "doing", "", false, "list doing tasks")
	listCmd.Flags().BoolVarP(&done, "done", "", false, "list done tasks")
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "list all tasks")

	rootCmd.AddCommand(listCmd)
}
