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
	alphaNumSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz_0123456789"
)

func main() {
	var length, batchSize int
	var format string
	flag.IntVar(&length, "length", 128, "input password length")
	flag.StringVar(&format, "format", "text", `input result type ("text", "json", "bytes", "hex", "base64")`)
	flag.IntVar(&batchSize, "batch", 1, "input batch size")
	flag.Parse()

	var result = make(chan []byte, length)
	defer close(result)
	go func() {
		for r := range result {
			fmt.Printf("%v\n", display(format, r))
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result <-generateBytes(length)
		}()
	}
	wg.Wait()
}

func generateBytes(length int) []byte {
	// generate password
	var wg sync.WaitGroup
	var result = make([]byte, length)
	result[0] = generateChar()
	result[length-1] = generateChar()
	for i := 1; i < length-1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for result[id] == result[id-1] || result[id] == result[id+1] {
				result[id] = generateChar()
			}
		}(i)
	}
	wg.Wait()
	return result
}

func generateSeed() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return time.Now().UnixNano() + int64(m.Alloc) + int64(m.Frees)
}

func generateChar() uint8 {
	// prepare seed
	seed := generateSeed()
	s := rand.NewSource(seed)
	r := rand.New(s)

	// select char
	selected := r.Intn(len(alphaNumSet))
	return alphaNumSet[selected]
}

func display(format string, result []byte) (display interface{}) {
	switch format {
	case "text":
		return string(result)
	case "hex":
		return hex.EncodeToString(result)
	case "base64":
		return base64.StdEncoding.EncodeToString(result)
	case "json":
		buff, _ := json.Marshal(
			map[string]interface{}{
				"length": len(result),
				"result": string(result),
			})
		return string(buff)
	case "bytes":
		return result
	default:
		return errors.New("unknown result format")
	}
}
