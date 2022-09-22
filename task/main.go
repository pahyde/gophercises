package main

import (
//    "fmt"
//    "log"
    "encoding/binary"
//    "github.com/boltdb/bolt"
//    "github.com/spf13/cobra"

    "task/cmd"
)

func main() {
    cmd.Execute()
}

/*
    fmt.Println("Hello task!")

    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("tasks"))
        if err != nil {
            return err
        }
        return nil
    })

    for i := 0; i < 10; i++ {
        err := db.Update(func(tx *bolt.Tx) error {
            b := tx.Bucket([]byte("tasks"))
            id, _ := b.NextSequence()
            err := b.Put(itob(id), []byte(fmt.Sprintf("task #%d", 9-i)))
            if err != nil {
                return err
            }
            return nil
        })
        if err != nil {
            log.Fatal(err)
        }
    }
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("tasks"))
        b.ForEach(func(k, v []byte) error {
		    fmt.Printf("key=%s, value=%s\n", k, v)
		    return nil
		})
	    return nil
    }) 
    if err != nil {
        log.Fatal(err)
    }

    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("tasks"))
        b.ForEach(func(k, v []byte) error {
		    fmt.Printf("key=%s, value=%s\n", k, v)
		    return nil
		})
	    return nil
    }) 
    if err != nil {
        log.Fatal(err)
    }

*/
func itob(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, v)
    return b
}
