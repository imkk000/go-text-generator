package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	alpha = iota
	number
)

const (
	alphaSet  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	numberSet = "0123456789"
)

func main() {
	var l int
	var f string
	flag.IntVar(&l, "length", 128, "input password length")
	flag.StringVar(&f, "format", "text", `input result type ("text", "json", "bytes", "hex", "base64")`)
	flag.Parse()

	// generate password
	var wg sync.WaitGroup
	var result = make([]byte, l)
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func(id int) {
			defer wg.Done()
			result[id] = generateChar()
		}(i)
	}
	wg.Wait()

	// check nearly duplicate
	for i := 1; i < l; i++ {
		if result[i] == result[i-1] {
			result[i] = generateChar()
			// update and recheck again
			i--
		}
	}

	// display result
	var display interface{}
	switch f {
	case "text":
		display = string(result)
	case "hex":
		display = hex.EncodeToString(result)
	case "base64":
		display = base64.StdEncoding.EncodeToString(result)
	case "json":
		buff, _ := json.Marshal(
			map[string]interface{}{
				"length": l,
				"result": string(result),
			})
		display = string(buff)
	case "bytes":
		display = result
	default:
		display = errors.New("unknown result format")
	}
	fmt.Printf("%v\n", display)
}

func generateSeed() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return time.Now().UnixNano() + int64(m.Alloc)
}

func generateChar() uint8 {
	// prepare seed
	seed := generateSeed()
	s := rand.NewSource(seed)
	r := rand.New(s)
	ii := r.Intn(2)

	// select char
	var charSet string
	switch ii {
	case alpha:
		charSet = alphaSet
	case number:
		charSet = numberSet
	}
	selected := r.Intn(len(charSet))
	return charSet[selected]
}
