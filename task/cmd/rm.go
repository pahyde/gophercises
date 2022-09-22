package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
  Use:   "rm",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("rm")
  },
}

func init() {
    rootCmd.AddCommand(rmCmd)
}
