package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
)

// Assumption: the alphabet can't be larger than 64 characters.
const AlphabetSize = 26
type Word struct {
    mask uint64
    occurrences [AlphabetSize]uint
    numChars uint
    mostUncommonChar byte
    text string
}

func (w Word) String() string {
    space := ""
    if len(w.text) <= 4 { space += "\t" }
    if len(w.text) <= 12 { space += "\t" }
    if len(w.text) <= 20 { space += "\t" }
    return fmt.Sprintf("\"%s\":%s%d characters\tmask %026b\toccurrences %v\tmost uncommon char %c", w.text, space, w.numChars, w.mask, w.occurrences, w.mostUncommonChar)
}

type CharacterData struct {
    occurrences [AlphabetSize]uint64
    frequencies [AlphabetSize]float64
}

var Alphabet = new(CharacterData)

var Dictionary []*Word

var dictionaryFile *string = flag.String("dict", "/usr/share/dict/words", "Dictionary file to use")

// character to index
func ctoi(char byte) uint {
    return uint(char - 'a')
}
// index to character
func itoc(index uint) byte {
    return byte(index + 'a')
}

func readWords() {
    file, err := os.Open(*dictionaryFile, os.O_RDONLY, 0000)
    defer file.Close()
    if err != nil {
        os.Exit(1)
    }
    reader := bufio.NewReader(file)

    // TODO: this should be done better
    charRegexp, err := regexp.Compile("^[a-z]+$")
    if err != nil {
        os.Exit(1)
    }

    for {
        if word, err := reader.ReadString('\n'); err == nil {
            word = strings.TrimSpace(word)
            if !charRegexp.MatchString(word) {
                continue
            }
            w := MakeWord(word)
            Dictionary = append(Dictionary, w)
        } else {
            break
        }
    }
}

func processAlphabetData() {
    // count characters
    for _, w := range Dictionary {
        for i := 0; i < len(w.text); i++ {
            index := ctoi(w.text[i])
            Alphabet.occurrences[index] += 1
        }
    }

    var total uint64 = 0
    for i := 0; i < AlphabetSize; i++ {
        total += Alphabet.occurrences[i]
    }

    // calculate frequencies
    for i := 0; i < AlphabetSize; i++ {
        Alphabet.frequencies[i] = float64(Alphabet.occurrences[i]) / float64(total)
    }
}


func MakeWord(text string) *Word {
    w := new(Word)
    w.text = text
    w.generateMaskAndOccurrences()
    return w
}

func (w *Word) generateMaskAndOccurrences() {
    for i := 0; i < len(w.text); i++ {
        index := ctoi(w.text[i])
        w.mask |= (1 << index)
        w.occurrences[index] += 1
        w.numChars += 1
    }
}

func (w *Word) determineMostUncommonCharacter() {
    var freq float64 = 1.0
    var char byte = 'a'

    for i := uint(0); i < AlphabetSize; i++ {
        if w.mask & (1 << i) != 0 {
            if Alphabet.frequencies[i] < freq {
                freq = Alphabet.frequencies[i]
                char = itoc(i)
            }
        }
    }

    w.mostUncommonChar = char
}

func main() {
    readWords()
    processAlphabetData()

    for _, w := range Dictionary {
        w.determineMostUncommonCharacter()
    }

    for _, w := range Dictionary {
        fmt.Printf("%+v\n", w)
    }

    fmt.Printf("%+v\n", Alphabet)
}

