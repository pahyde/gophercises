package main

import (
    "fmt"
    "strings"
    "os"
    "flag"
    "encoding/csv"
    "time"
)

type Problem struct {
    q string
    a string
}

type Answer struct {
    isCorrect bool
}

func main() {

    filename := *flag.String(
        "csv", 
        "problems.csv", 
        "a csv file in the format (question, answer). defaults to problems.csv",
    )
    seconds := *flag.Int(
        "time", 
        16,
        "the time limit for the quiz (in seconds). defaults to 16s",
    )
    flag.Parse()

    file, err := os.Open(filename)
    if err != nil {
        exit(fmt.Sprintf("Failed to open file: %s!", filename))
    }

    r := csv.NewReader(file)
    lines, err := r.ReadAll()
    if err != nil {
        exit(fmt.Sprintf("Something went wrong reading %s.", filename))
    }

    problems, err := parseProblems(lines)
    if err != nil {
        exit(err.Error())
    }

    ansCh  := make(chan Answer)
    quitCh := make(chan int)
    
    // handles state changes (e.g. submitted answers and time limit reached)
    go handleState(ansCh, quitCh, len(problems))

    // start timer
    // note: time.NewTimer might be a better choice here
    go func() {
        time.Sleep(time.Duration(seconds) * time.Second)
        quitCh <- 0
    }()

    // start quiz
    fmt.Println("Ready, set, Go! (pun intended)")
    for _, p := range problems {
        answer, err := getUserAnswer(p)
        if err != nil {
            exit("Something went wrong reading your answer. Please try again.")
        }
        ansCh <- Answer{isCorrect: answer == p.a}
    }
    fmt.Print("\nNice job!")
    quitCh <- 0
}

func handleState(ansCh chan Answer, quitCh chan int, total int) {
    correct  := 0
    answered := 0
    for {
        select {
        case a := <-ansCh:
            answered++
            if a.isCorrect {
                correct++
            }
        case <-quitCh:
            displayResults(correct, answered, total)
            return
        }
    }
}

func displayResults(correct, answered, total int) {
    score := 0.0
    if answered > 0 {
        score = float64(correct) / float64(answered) * 100
    }
    template := "\n\nYou answered %d out of %d questions and got %d/%d (%.2f%%) correct\n\n"
    fmt.Printf(template, correct, total, correct, answered, score)
    os.Exit(0)
}

func getUserAnswer(p Problem) (string, error) {
    fmt.Printf("%s = ", p.q)
    var answer string
    _, err := fmt.Scanln(&answer)
    if err != nil {
        return answer, err
    }
    return answer, nil
}

func parseProblems(lines [][]string) ([]Problem, error) {
    problems := make([]Problem, len(lines))
    for i, line := range lines {
        if len(line) != 2 {
            return nil, fmt.Errorf("quiz question %d is formatted incorrectly. Must contain a question and answer column only.", i)
        }
        p := Problem{
            q: line[0],
            a: strings.TrimSpace(line[1]),
        }
        problems[i] = p
    }
    return problems, nil
}

func exit(m string) {
    fmt.Println(m)
    os.Exit(1)
}
