package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	RuneSets            = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_+=!-[]\\/()"
	DefaultRuneSelector = `\w+`
)

func main() {
	var runeSelector string
	var pwdLen, batchSize int
	flag.StringVar(&runeSelector, "selector", DefaultRuneSelector, "input custom rune selector")
	flag.IntVar(&pwdLen, "length", 128, "input each password length (1-1000)")
	flag.IntVar(&batchSize, "batch", 1, "input batch size (1-10000)")
	flag.Parse()

	if runeSelector == "" {
		log.Fatal("invalid custom rune selector")
	}
	if pwdLen <= 0 || pwdLen > 1000 {
		log.Fatal("invalid password length")
	}
	if batchSize <= 0 || batchSize > 10000 {
		log.Fatal("invalid batch size")
	}

	re, err := regexp.Compile(runeSelector)
	if err != nil {
		log.Fatalf("new regex %v", err)
	}
	runes := re.FindString(RuneSets)
	if len(runes) == 0 {
		log.Fatal("parse rune sets error")
	}

	ch := make(chan string, batchSize)
	go generatePwd(ch, batchSize, pwdLen, []rune(runes))

	for s := range ch {
		fmt.Println(s)
	}
}

func generatePwd(sender chan<- string, s, l int, runes []rune) {
	defer close(sender)

	mut := new(sync.Mutex)
	wgOuter := new(sync.WaitGroup)
	wgOuter.Add(s)
	for s > 0 {
		s--

		go func() {
			defer wgOuter.Done()

			wgInner := new(sync.WaitGroup)
			wgInner.Add(l)
			var b strings.Builder
			for i := 0; i < l; i++ {
				go func() {
					defer wgInner.Done()

					src := rand.NewSource(seed())
					j := src.Int63() % int64(len(runes))

					mut.Lock()
					b.WriteRune(runes[j])
					mut.Unlock()
				}()
			}
			wgInner.Wait()

			sender <- b.String()
		}()
	}
	wgOuter.Wait()
}

func seed() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return time.Now().UnixNano() + int64(m.Alloc) + int64(m.Frees)
}
