package cmd

import (
    "fmt"
    "log"
    "github.com/spf13/cobra"

    "task/tasklist"
)

var completedCmd = &cobra.Command{
  Use:   "completed",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      l, err := tasklist.Open()
      if err != nil {
          log.Fatal(err)
      }
      defer func() {
          if err := l.Close(); err != nil {
              log.Fatal(err)
          }
      }()
      completed, err := l.Completed()
      if err != nil {
          log.Fatal(err)
      }
      fmt.Println("You have finished the following tasks today:")
      for _, tsk := range completed {
          fmt.Printf("- %s\n", tsk.Name)
      }
  },
}

func init() {
    rootCmd.AddCommand(completedCmd)
}
