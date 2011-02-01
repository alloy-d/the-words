package dictionary

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
)

// character to index
func CtoI(char byte) uint {
    return uint(char - 'a')
}
// index to character
func ItoC(index uint) byte {
    return byte(index + 'a')
}

// Assumption: the alphabet can't be larger than 64 characters.
const AlphabetSize = 26
type LetterBag struct {
    Mask uint64
    Occurrences [AlphabetSize]uint
    NumLetters uint
    MostUncommonLetter byte
}

type Word struct {
    Text string
    LetterBag
}

func NewWord(text string) *Word {
    w := new(Word)
    w.Text = text
    w.generateMaskAndOccurrences(w.Text)
    return w
}

func (b *LetterBag) generateMaskAndOccurrences(text string) {
    for i := 0; i < len(text); i++ {
        index := CtoI(text[i])
        b.Mask |= (1 << index)
        b.Occurrences[index] += 1
        b.NumLetters += 1
    }
}

func (b *LetterBag) determineMostUncommonCharacter(chars *CharacterData) {
    var freq float64 = 1.0
    var char byte = 'a'

    for i := uint(0); i < AlphabetSize; i++ {
        if b.Mask & (1 << i) != 0 {
            if chars.Frequencies[i] < freq {
                freq = chars.Frequencies[i]
                char = ItoC(i)
            }
        }
    }

    b.MostUncommonLetter = char
}

func (b LetterBag) String() string {
    return fmt.Sprintf("%d letters\tmask %026b\toccurrences %v\tmost uncommon letter %c", b.NumLetters, b.Mask, b.Occurrences, b.MostUncommonLetter)
}

func (w Word) String() string {
    space := ""
    if len(w.Text) <= 4 { space += "\t" }
    if len(w.Text) <= 12 { space += "\t" }
    if len(w.Text) <= 20 { space += "\t" }
    return fmt.Sprintf("\"%s\":%s%v", w.Text, space, w.LetterBag)
}


type CharacterData struct {
    Occurrences [AlphabetSize]uint64
    Frequencies [AlphabetSize]float64
}

func (c CharacterData) String() string {
    s := "Character Data: --------------------------------------------\n"
    for i := uint(0); i < AlphabetSize; i++ {
        s += fmt.Sprintf("  '%c' occurred %10d times with a frequency of %0.4f\n", ItoC(i), c.Occurrences[i], c.Frequencies[i])
    }
    s += "------------------------------------------------------------\n"

    return s
}


type Dictionary struct {
    Words []*Word
    *CharacterData
}

func New(fileName string) *Dictionary {
    d := new(Dictionary)
    d.CharacterData = new(CharacterData)
    d.readWords(fileName)
    d.processAlphabetData()
    return d
}

func (d *Dictionary) readWords(fileName string) {
    file, err := os.Open(fileName, os.O_RDONLY, 0000)
    defer file.Close()
    if err != nil {
        panic(err)
    }
    reader := bufio.NewReader(file)

    // TODO: this should be done better
    charRegexp, err := regexp.Compile("^[a-z]+$")
    if err != nil {
        panic(err)
    }

    for {
        if word, err := reader.ReadString('\n'); err == nil {
            word = strings.TrimSpace(word)
            if !charRegexp.MatchString(word) {
                continue
            }
            w := NewWord(word)
            d.Words = append(d.Words, w)
        } else {
            break
        }
    }
}

func (d *Dictionary) processAlphabetData() {
    // count characters
    for _, w := range d.Words {
        for i := 0; i < len(w.Text); i++ {
            index := CtoI(w.Text[i])
            d.Occurrences[index] += 1
        }
    }

    var total uint64 = 0
    for i := 0; i < AlphabetSize; i++ {
        total += d.Occurrences[i]
    }

    // calculate frequencies
    for i := 0; i < AlphabetSize; i++ {
        d.Frequencies[i] = float64(d.Occurrences[i]) / float64(total)
    }

    for _, w := range d.Words {
        w.determineMostUncommonCharacter(d.CharacterData)
    }
}

