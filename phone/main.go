package main

import (
    "fmt"
    "log"
    "os"
    "unicode"
    "io"
    "encoding/csv"
    "database/sql"

    _ "github.com/lib/pq"
)

type Person struct {
    FirstName   string
    LastName    string
    PhoneNumber string
}

func main() {
    const (
        host string = "localhost"
        port        = "5432"
        user        = "normalize"
        password    = "0112358"
        dbname      = "phone"
    )

    db, err := postgresDB(host, port, user, password, dbname)
    defer func() {
        if err := db.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    // Init numbers table to rows in numbers.csv
    db.Exec("DELETE FROM numbers")
    
    f, err := os.Open("numbers.csv")
    if err != nil {
        log.Fatal(err)
    }
    people, err := getPeopleFromCsv(f)
    if err != nil {
        log.Fatal(err)
    }
    if err := insertTableRows(db, people); err != nil {
        log.Fatal(err)
    }

    // get (id, phone number) for every row in table
    rows, err := db.Query("SELECT id, phone_number FROM numbers")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // Normalize phone numbers
    // delete duplicates, otherwise update w/ normalized phone number
    seen := make(map[string]bool)
    for rows.Next() {
        var id           int
        var phoneNumber  string
        if err := rows.Scan(&id, &phoneNumber); err != nil {
            log.Fatal(err)
        }
        n := normalize(phoneNumber)
        if seen[n] {
            db.Exec("DELETE FROM numbers WHERE id = $1", id)
        } else {
            db.Exec("UPDATE numbers SET phone_number = $1 WHERE id = $2", n, id)
        }
        seen[n] = true
    }
}

func normalize(phoneNumber string) string {
    digits := make([]rune, 0)
    for _, c := range phoneNumber {
        if unicode.IsDigit(c) {
            digits = append(digits, c)
        }
    }
    return string(digits)
}

func insertTableRows(db *sql.DB, people []Person) error {
    qstr := "INSERT INTO numbers (first_name, last_name, phone_number) VALUES ($1, $2, $3)"
    for _, p := range people {
        _, err := db.Exec(qstr, p.FirstName, p.LastName, p.PhoneNumber)
        if err != nil {
            return err
        }
    }
    return nil
}

func getPeopleFromCsv(r io.Reader) ([]Person, error) {
    records, err := csv.NewReader(r).ReadAll()
    if err != nil {
        return nil, err
    }
    var people []Person
    for _, rec := range records {
        if len(rec) != 3 {
            return nil, fmt.Errorf("Bad csv format.")
        }
        p := Person{rec[0], rec[1], rec[2]}
        people = append(people, p)
    }
    return people, nil
}

func postgresDB(host, port, user, password, dbname string) (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
        host, port, user, password, dbname,
    )
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    fmt.Println("Connected to db successfully!")
    return db, nil
}

