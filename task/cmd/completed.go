package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
  Use:   "completed",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("completed")
  },
}

func init() {
    rootCmd.AddCommand(completedCmd)
}
