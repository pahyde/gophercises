package tasklist

import (
    "github.com/boltdb/bolt"
)

type TaskList struct {
    db *bolt.DB
}

type Task struct {
    Name  string
    ta    string  //time: RFC3339
    td    string  //time: RFC3339
}

func (tsk Task) TimeAdd() (time.Time, error) {
    return time.Parse(time.RFC3339, tsk.ta)
}

func (tsk Task) TimeDone() (time.Time, error) {
    return time.Parse(time.RFC3339, tsk.td)
}

func Open() (*TaskList, error) {
    db, err := bolt.Open("tasks.db", 0600, nil)
    if err != nil {
        return nil, err
    }
    err := initBuckets(db, "todo", "completed")
    if err != nil {
        return nil, err
    }
    return &TaskList{db}
}

func (l *TaskList) Close() error {
    return l.db.Close()
}

// (t, task) -> todo bucket
func (l *TaskList) Add(string taskName) error {
    return l.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("todo"))
        t := time.Now().Format(time.RFC3339)
        tsk := Task{taskName, t, ""}
        return putTaskToBucket(t, tsk, b)
    })
}

// todo bucket -> (t, task) -> completed bucket
func (l *TaskList) Done(n int) (Task, error) {
    var done string
    err := l.db.Update(func(tx *bolt.Tx) error {
        b1 := tx.Bucket([]byte("todo"))
        b2 := tx.Bucket([]byte("completed"))
        // find and remove tsk #n
        tsk, err := findAndRemove(n, b1)
        if err != nil {
            return err
        }
        // update time done
        tsk.td = time.Now().Format(time.RFC3339)
        // set return value
        done = tsk
        // put to completed
        return putTaskToBucket(tsk.td, tsk, b2)
    })
    return done, err
}

// todo bucket -> (t, task) -> rm
func (l *TaskList) Rm(n int) (Task, error) {
    var removed string
    err := l.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("todo"))
        // find and remove tsk #n
        tsk, err := findAndRemove(n, b)
        if err != nil {
            return err
        }
        // set return value
        removed = tsk
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
        })
    })
    return tasks, err
}

// TODO: iterate all kv, delete yesterdays tasks and collect todays
//       then return todays asks
func (l *TaskList) Completed(n int) ([]Task, error) {
    tasks := make([]Task, 0)
    err := l.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("completed"))
        // get tasks from completed bucket
        return b.ForEach(func(k, v []byte) error {
            var tsk Task
            if err := json.Unmarshal(v, &tsk); err != nil {
                return err
            }
            tasks = append(tasks, tsk)
        })
    })
    return tasks, err
}

// put (t, tsk) to bucket b
func putTaskToBucket(t string, tsk Task, b *bolt.Bucket) error {
        dat, err := json.Marshal(&tsk) 
        if err != nil {
            return err
        }
        return b.Put([]byte(t), tsk)
}

// deletes nth task and returns as Task
func findAndRemove(n int, b *bolt.Bucket) (Task, error) {
    c := b.Cursor()
    i := 1
    // get todo (t, task) at index idx 
    for t, tsk := c.First(); t != nil; t, tsk := c.Next() {
        if (i < n) {
            i++
            continue
        }
        // delete and return (t, tsk) from todo
        found, err := json.Marshal(&tsk)
        if err != nil {
            return "", err
        }
        if err := b1.Delete(t); err != nil {
            return "", err
        }
        return found, nil
    }
    return "", fmt.Errorf("argument idx: %d out of bounds", idx)
}

func initBuckets(db *bolt.DB, buckets ...string) error {
    return db.Update(func(tx *bolt.Tx) error {
        for _, bkt := range buckets {
            _, err := tx.CreateBucketIfNotExsits([]byte(bkt))
            if err != nil {
                return err
            }
        }
    })
}
