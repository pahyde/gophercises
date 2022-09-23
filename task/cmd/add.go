package cmd

import (
    "fmt"
    "log"
    "strings"
    "github.com/spf13/cobra"

    "task/tasklist"
)

var addCmd = &cobra.Command{
  Use:   "add",
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
      taskName := strings.Join(args, " ")
      l.Add(taskName)
      fmt.Printf("Added \"%s\" to your task list.\n", taskName)
  },
}

func init() {
    rootCmd.AddCommand(addCmd)
}
