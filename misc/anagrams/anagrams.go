package anagrams

import (
    "alloy-d/dictionary"
    "os"
    "regexp"
)

type Pool dictionary.LetterBag
type Anagram struct {
    Text string
    pool Pool
}

var charRegexp *regexp.Regexp
func init() {
    var err os.Error
    charRegexp, err = regexp.Compile("[a-z]")
    if err != nil {
        panic(err)
    }
}

func (p *Pool) AddWord(word string) {
    for i := 0; i < len(word); i++ {
        if charRegexp.Match([]byte{word[i]}) {
            p.AddChar(word[i])
        }
    }
}

func (p *Pool) AddChar(c byte) {
    index := dictionary.CtoI(c)
    p.Mask |= (1 << index)
    p.Occurrences[index] += 1
    p.NumLetters += 1
}

func (p Pool) ContainsWord(w *dictionary.Word) bool {
    if p.Mask | w.Mask != p.Mask { return false }
    if p.NumLetters < w.NumLetters { return false }
    for i := 0; i < dictionary.AlphabetSize; i++ {
        if p.Occurrences[i] < w.Occurrences[i] {
            return false
        }
    }
    return true
}

func (p *Pool) FindAnagrams(dict *dictionary.Dictionary, out chan<- string) {
    a := new(Anagram)
    a.pool = *p

    q := StartQueue(dict, out)
    q <- *a
}

func (a Anagram) FindAnagrams(dict *dictionary.Dictionary, queue chan<- Anagram, out chan<- string) {
    if a.pool.NumLetters == 0 { out <- a.Text }
    for _, w := range dict.Words {
        if w.NumLetters < 3 { continue }
        if a.pool.ContainsWord(w) {
            queue <- a.AddWord(w)
        }
    }
}

func (a Anagram) AddWord(w *dictionary.Word) Anagram {
    for i := 0; i < dictionary.AlphabetSize; i++ {
        a.pool.Occurrences[i] -= w.Occurrences[i]
    }
    a.pool.NumLetters -= w.NumLetters
    if a.Text != "" { a.Text += " " }
    a.Text += w.Text
    return a
}

