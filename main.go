package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	alphaSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numberSet  = "0123456789"
	specialSet = "!@#$%*="
)

const (
	alpha = iota
	number
	special
)

func main() {
	var l int
	flag.IntVar(&l, "length", 128, "input password length")

	// prepare seed
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(runtime.NumGoroutine())

	seed := time.Now().Unix() + int64(m.Alloc)
	s := rand.NewSource(seed)
	r := rand.New(s)

	var wg sync.WaitGroup
	var result = make([]byte, l)
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func(id int) {
			defer wg.Done()
			var charSet string
			ii := r.Intn(2)
			switch ii {
			case alpha:
				charSet = alphaSet
			case number:
				charSet = numberSet
			}
			ij := r.Intn(len(charSet))
			result[id] = charSet[ij]
		}(i)
	}
	wg.Wait()
	fmt.Println(string(result))
}
