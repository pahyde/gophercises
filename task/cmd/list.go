package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
  Use:   "list",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("list")
  },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
