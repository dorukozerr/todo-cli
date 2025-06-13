package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple todo CLI",
	Long:  "A command-line todo application with group management and priority levels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Todo CLI - Use 'todo --help' for available commands")
	},
}

func init() {
	RootCmd.SilenceUsage = true
	RootCmd.SilenceErrors = true

	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(completeCmd)
	RootCmd.AddCommand(incompleteCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(groupCmd)
}
