package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Verbose bool

	rootCmd = &cobra.Command{
		Use:     "task",
		Short:   "Manage your daily tasks in cli",
		Long:    "`task` is a cli tool helping you manage your daily tasks, todos, missions in a single place",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("task")
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "turn on verbose mode")

	rootCmd.AddGroup(&cobra.Group{ID: "task-life-cycle", Title: "Task life cycle commands"})
	rootCmd.AddGroup(&cobra.Group{ID: "general", Title: "General commands"})
}

func Execute() error {
	return rootCmd.Execute()
}
