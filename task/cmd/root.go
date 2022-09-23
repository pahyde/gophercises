package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
    "os"
)


/*
commands:
    add:       (t, todo) -> todos
    rm:        todos -> (t, todo)                [linear search]
    do:        todos -> (t, todo) -> completed   [linear search]
    list:      todos -> stdout
        overdue todos: red(todo <- (added [yesterday | n days ago]))
    completed: completed -> stdout
        removes previous days completed todos (up to 12am)


    k: time: RFC3339
    v: seialized todo ([]byte)
*/

var rootCmd = &cobra.Command{
  Use:   "task",
  Short: "task is a CLI for managing your daily TODOs.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(10)
    }
}
