package cmd

import (
    "log"
    "strings"
    "github.com/boltdb/bolt"
    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
  Use:   "add",
  Short: "task is a CLI for managing your daily TODOs.",
  Run: func(cmd *cobra.Command, args []string) {
      db, err := bolt.Open("tasks.db", 0600, nil)
      if err != nil {
          log.Fatal(err)
      }
      defer db.Close()
      initBuckets(db, "todos", "completed")
      db.Update(func(tx *bolt.Tx) error {
          b := tx.Bucket([]byte("todos"))
          t := time.Now().Format(time.RFC3339)
          todo := strings.Join(args, " ")
          b.Put(t, todo)
      })
  },
}

func init() {
    rootCmd.AddCommand(addCmd)
}
