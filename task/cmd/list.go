package cmd

import (
    "time"
    "fmt"
    "log"
    "github.com/spf13/cobra"

    "task/tasklist"
)

var listCmd = &cobra.Command{
  Use:   "list",
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
      list, err := l.List()
      if err != nil {
          log.Fatal(err)
      }

      fmt.Println("You have the following tasks:")
      for i, tsk := range list {
          // late message for overdue tasks, empty if added today
          m, err := taskMessage(tsk)
          if err != nil {
              log.Fatal(err)
          }
          fmt.Printf("%d. %s%s\n", i+1, tsk.Name, m)
      }
  },
}

func taskMessage(tsk tasklist.Task) (string, error) {
      n, err := daysLate(tsk)
      if err != nil {
          return "", err
      }
      switch {
      case n > 1:
          return fmt.Sprintf(" (added %d days ago)", n), nil
      case n == 1: 
          return fmt.Sprintf(" (added yesterday)"), nil
      default:
          return "", nil
      }
}

func daysLate(tsk tasklist.Task) (int, error) {
    now := time.Now()
    y, m, d  := now.Date()
    midnight := time.Date(
        y, m, d,
        0, 0, 0, 0,
        now.Location(),
    )
    t, err := tsk.TimeAdded()
    if err != nil {
        return 0, err
    }
    if t.After(midnight) {
        return 0, nil
    }
    days := int(midnight.Sub(t).Hours()) / 24 + 1
    return days, nil
}

func init() {
    rootCmd.AddCommand(listCmd)
}
