package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
)

type intermediate struct {
    cur string
    chars []int
}

var wordsFile *string = flag.String("dict", "/usr/share/dict/words", "Dictionary file to use")
var words map[string]bool = make(map[string]bool)

// Character to index.
func ctoi(char byte) int {
    return int(char - 'a')
}
// Index to character.
func itoc(index int) byte {
    return byte(index + 'a')
}

func hasMoreLetters(d *intermediate) bool {
    for i := 0; i < len(d.chars); i++ {
        if d.chars[i] > 0 {
            return true
        }
    }
    return false
}

func anagramDispatcher(out chan<- string, d *intermediate) {
    numRoutines := 0
    in := make(chan string)
    for i := 0; i < 26; i++ {
        if d.chars[i] > 0 {
            n := new(intermediate)
            n.cur = string(itoc(i))
            n.chars = append([]int{}, d.chars...)
            n.chars[i] -= 1
            numRoutines++
            go findAnagrams(in, n)
        }
    }

    for word := range in {
        if word == "" {
            numRoutines--
        }
        if numRoutines == 0 {
            break;
        }

        if d.cur == "" {
            out <- d.cur + word
        } else {
            out <- d.cur + " " + word
        }
    }
    close(in)
}

func findAnagrams(out chan<- string, d *intermediate) {
    if words[d.cur] {
        if (hasMoreLetters(d)) {
            anagramDispatcher(out, d)
        } else {
            out <- d.cur
        }
    }
    for i := 0; i < 26; i++ {
        if d.chars[i] > 0 {
            n := new(intermediate)
            n.cur = d.cur + string(itoc(i))
            n.chars = append([]int{}, d.chars...)
            n.chars[i] -= 1
            findAnagrams(out, n)
        }
    }
    out <- ""
}

func readWords() {
    file, err := os.Open(*wordsFile, os.O_RDONLY, 0000)
    if err != nil {
        os.Exit(1)
    }
    reader := bufio.NewReader(file)

    for {
        if word, err := reader.ReadString('\n'); err == nil {
            words[strings.TrimSpace(word)] = true
        } else {
            break
        }
    }
}

func main() {
    flag.Parse()
    readWords()
    words[""] = false

    var chars [26]int
    charRegexp, err := regexp.Compile("[a-z]")
    if err != nil {
        os.Exit(1)
    }

    for i := 0; i < flag.NArg(); i++ {
        word := strings.ToLower(flag.Arg(i))
        for j := 0; j < len(word); j++ {
            if charRegexp.Match([]byte{word[j]}) {
                chars[ctoi(word[j])] += 1
            }
        }
    }

    /*
    for i := 0; i < len(chars); i++ {
        fmt.Printf("%c: %d\n", itoc(i), chars[i])
    }
    */

    results := make(chan string)
    go func() {
        anagramDispatcher(results, &intermediate{cur: "", chars: chars[0:26]})
    }()
    for word := range results {
        fmt.Println(word)
    }
}

