package main

import (
    "alloy-d/anagrams"
    "alloy-d/dictionary"
    "flag"
    "fmt"
)

var dictionaryFile *string = flag.String("dict", "/usr/share/dict/words", "Dictionary file to use")
var printWords *bool = flag.Bool("print-words", false, "Whether to print all word data")
var printChars *bool = flag.Bool("print-chars", false, "Whether to print alphabet data")

func main() {
    flag.Parse()
    d := dictionary.New(*dictionaryFile)

    if *printWords {
        for _, w := range d.Words {
            fmt.Printf("%+v\n", w)
        }
    }

    if *printChars {
        fmt.Printf("%+v", d)
    }

    if flag.NArg() > 0 {
        in := make(chan string)
        p := new(anagrams.Pool)
        p.AddWord(flag.Arg(0))
        go p.FindAnagrams(d, in)

        for {
            select {
                case a := <-in:
                    fmt.Println(a)
            }
        }
    }
}

