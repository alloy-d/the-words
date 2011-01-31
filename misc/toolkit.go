package main

import (
    "./dictionary"
    "flag"
    "fmt"
)

var dictionaryFile *string = flag.String("dict", "/usr/share/dict/words", "Dictionary file to use")
var printWords *bool = flag.Bool("print-words", false, "Whether to print all word data")
var printChars *bool = flag.Bool("print-chars", true, "Whether to print alphabet data")

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
}

