package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
  Use:   "do",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("do")
  },
}

func init() {
    rootCmd.AddCommand(doCmd)
}
