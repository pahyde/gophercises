package main

import (
    "fmt"
    "log"
    "os"
    "io"
    "flag"
    "encoding/csv"
)

func main() {
    total   := 0
    correct := 0

    filename := flag.String(
        "csv", 
        "problems.csv", 
        "a csv file in the format (question, answer). defaults to problems.csv",
    )
    flag.Parse()

    file, err := os.Open(*filename)
    if err != nil {
        log.Fatal(err)
    }

    r := csv.NewReader(file)
    fmt.Println("Ready, set, Go! (pun intended)")
    for {
        entry, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        if len(entry) != 2 {
            err := fmt.Errorf("quiz question %d is formatted incorrectly. Must contain a question and answer column only.", len(entry))
            log.Fatal(err)
        }

        question   := entry[0]
        answer     := entry[1]
        userAnswer := ""

        fmt.Printf("%s : ", question)
        _, err = fmt.Scanln(&userAnswer)
        if err != nil {
            log.Fatal(err)
        }
        if userAnswer == answer {
            correct++
        }
        total++
    }

    fmt.Printf("Nice, you answered %d/%d correctly\n", correct, total)
}
