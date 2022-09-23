package tasklist

import (
    "fmt"
    "time"
    "bytes"
    "encoding/json"
    "github.com/boltdb/bolt"
)

type TaskList struct {
    db *bolt.DB
}

type Task struct {
    Name  string
    Ta    string  //time added: RFC3339
    Td    string  //time done:  RFC3339
}

func (tsk Task) TimeAdded() (time.Time, error) {
    return time.Parse(time.RFC3339, tsk.Ta)
}

func (tsk Task) TimeDone() (time.Time, error) {
    return time.Parse(time.RFC3339, tsk.Td)
}

func Open() (*TaskList, error) {
    db, err := bolt.Open("tasks.db", 0600, nil)
    if err != nil {
        return nil, err
    }
    err = initBuckets(db, "todo", "completed")
    if err != nil {
        return nil, err
    }
    return &TaskList{db}, nil
}

func (l *TaskList) Close() error {
    return l.db.Close()
}

// (t, task) -> todo bucket
func (l *TaskList) Add(taskName string) error {
    return l.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("todo"))
        t := time.Now().Format(time.RFC3339)
        tsk := Task{taskName, t, ""}
        return putTaskToBucket(t, tsk, b)
    })
}

// todo bucket -> (t, task) -> completed bucket
func (l *TaskList) Done(n int) (Task, error) {
    var done Task
    err := l.db.Update(func(tx *bolt.Tx) error {
        b1 := tx.Bucket([]byte("todo"))
        b2 := tx.Bucket([]byte("completed"))
        // find and remove tsk #n
        tsk, err := findAndRemove(n, b1)
        if err != nil {
            return err
        }
        // update time done
        tsk.Td = time.Now().Format(time.RFC3339)
        // set return value
        done = tsk
        // put to completed
        return putTaskToBucket(tsk.Td, tsk, b2)
    })
    return done, err
}

// todo bucket -> (t, task) -> rm
func (l *TaskList) Rm(n int) (Task, error) {
    var removed Task
    err := l.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("todo"))
        // find and remove tsk #n
        tsk, err := findAndRemove(n, b)
        if err != nil {
            return err
        }
        // set return value
        removed = tsk
        return nil
    })
    return removed, err
}

func (l *TaskList) List() ([]Task, error) {
    tasks := make([]Task, 0)
    err := l.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("todo"))
        // get tasks from todo bucket
        return b.ForEach(func(k, v []byte) error {
            var tsk Task
            if err := json.Unmarshal(v, &tsk); err != nil {
                return err
            }
            tasks = append(tasks, tsk)
            return nil
        })
    })
    return tasks, err
}

// iterate all kv, delete yesterdays tasks and collect todays
// then return todays tasks
func (l *TaskList) Completed() ([]Task, error) {
    tasks := make([]Task, 0)
    err := l.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("completed"))
        m := []byte(midnight().Format(time.RFC3339))
        c := b.Cursor()
        // delete yesterday or later completed tasks
        for t, _ := c.First(); t != nil && bytes.Compare(t, m) < 0; t, _ = c.Next() {
            if err := b.Delete(t); err != nil {
                return err
            }
        }
        // collect todays completed tasks
        for t, dat := c.Seek(m); t != nil; t, dat = c.Next() {
            var tsk Task
            if err := json.Unmarshal(dat, &tsk); err != nil {
                return err
            }
            tasks = append(tasks, tsk)
        }
        return nil
    })
    return tasks, err
}

func midnight() time.Time {
    now     := time.Now()
    y, m, d := now.Date()
    return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

// put (t, tsk) to bucket b
func putTaskToBucket(t string, tsk Task, b *bolt.Bucket) error {
        dat, err := json.Marshal(&tsk) 
        if err != nil {
            return err
        }
        return b.Put([]byte(t), dat)
}

// deletes nth task and returns as Task
func findAndRemove(n int, b *bolt.Bucket) (Task, error) {
    c := b.Cursor()
    i := 1
    // get todo (t, task) at index idx 
    for t, tsk := c.First(); t != nil; t, tsk = c.Next() {
        if (i < n) {
            i++
            continue
        }
        // delete and return (t, tsk) from todo
        var found Task
        if err := json.Unmarshal(tsk, &found); err != nil {
            return Task{}, err
        }
        if err := b.Delete(t); err != nil {
            return Task{}, err
        }
        return found, nil
    }
    return Task{}, fmt.Errorf("argument idx: %d out of bounds", n)
}

func initBuckets(db *bolt.DB, buckets ...string) error {
    return db.Update(func(tx *bolt.Tx) error {
        for _, b := range buckets {
            _, err := tx.CreateBucketIfNotExists([]byte(b))
            if err != nil {
                return err
            }
        }
        return nil
    })
}
