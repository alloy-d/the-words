package anagrams

import (
    "alloy-d/dictionary"
    "container/heap"
    "container/vector"
    "runtime"
)

type queue struct {
    vector.Vector
}

func (q queue) Less(i, j int) bool {
    return q.At(i).(Anagram).pool.NumLetters < q.At(j).(Anagram).pool.NumLetters
}

func StartQueue(dict *dictionary.Dictionary, out chan<- string) chan<- Anagram {
    runtime.GOMAXPROCS(3)
    c := make(chan Anagram)

    go func () {
        q := new(queue)
        heap.Init(q)

        for {
            /*
            if q.Len() > 0 {
                a := heap.Pop(q)
                go a.(Anagram).FindAnagrams(dict, c, out)
            }
            */

            select {
                case a := <-c:
                    heap.Push(q, a)
                default:
                    for i := 0; i < 100 && q.Len() > 0; i++ {
                        a := heap.Pop(q)
                        go a.(Anagram).FindAnagrams(dict, c, out)
                    }
            }
        }
    }()

    return c
}

