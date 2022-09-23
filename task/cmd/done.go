package cmd

import (
    "fmt"
    "log"
    "strconv"
    "github.com/spf13/cobra"

    "task/tasklist"
)

var doneCmd = &cobra.Command{
  Use:   "done",
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
      num, err := strconv.Atoi(args[0])
      if err != nil {
          log.Fatal(err)
      }
      tsk, err := l.Done(num)
      fmt.Printf("You have completed the \"%s\" task.\n", tsk.Name)
  },
}

func init() {
    rootCmd.AddCommand(doneCmd)
}
