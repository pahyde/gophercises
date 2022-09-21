package main

import (
    "strings"
    "unicode"
)

func main() {

    // TODO: add some test case

}

func camelcase(s string) {
    if s == "" {
        return 0
    }
    count := 1
    for _, c := range s {
        if c < 'a' {
            count++
        }
    }
    return count
}


func caesarCipher(s string, k int32) string {
    var b strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) {
            from := (r | (1 << 5)) - 97
            to   := (from + k) % 26
            diff := to - from
            
            r += diff
        }
        b.WriteRune(r)
    }
    return b.String()
}

