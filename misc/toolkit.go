package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
)

// Assumption: the alphabet can't be larger than 64 characters.
type Word struct {
    mask uint64
    occurrences []int
    freq float64
    text string
}

var Dictionary []*Word

var dictionaryFile *string = flag.String("dict", "/usr/share/dict/words", "Dictionary file to use")

// character to index
func ctoi(char byte) int {
    return int(char - 'a')
}
// index to character
func itoc(index int) byte {
    return byte(index + 'a')
}

func readWords() {
    file, err := os.Open(*dictionaryFile, os.O_RDONLY, 0000)
    defer file.Close()
    if err != nil {
        os.Exit(1)
    }
    reader := bufio.NewReader(file)

    for {
        if word, err := reader.ReadString('\n'); err == nil {
            w := new(Word)
            w.text = strings.TrimSpace(word)
            Dictionary = append(Dictionary, w)
        } else {
            break
        }
    }
}

func generateMasks() {
    for _, w := range Dictionary {
        for i := 0; i < len(w.text); i++ {

        }
    }
}

func main() {
    readWords()

    for _, w := range Dictionary {
        fmt.Println(w.text)
    }
}
