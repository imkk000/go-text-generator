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
	RuneSets            = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	DefaultRuneSelector = `[\w]+`
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

	regex, err := regexp.Compile(runeSelector)
	if err != nil {
		log.Fatalf("new regex %v", err)
	}
	runeSets := regex.FindAllString(RuneSets, -1)
	if len(runeSets) == 0 {
		log.Fatal("parse rune sets error")
	}
	runes := strings.Join(runeSets, "")
	if len(runes) == 0 {
		log.Fatal("parse rune sets error")
	}

	ch := make(chan string, batchSize)
	go generatePwd(ch, batchSize, pwdLen, []byte(runes))

	for s := range ch {
		fmt.Println(s)
	}
}

func generatePwd(sender chan<- string, batchSize, pwdLen int, runes []byte) {
	defer close(sender)

	runeLen := len(runes)

	wgOuter := new(sync.WaitGroup)
	wgOuter.Add(batchSize)
	for batchSize > 0 {
		batchSize--

		go func() {
			defer wgOuter.Done()

			j := getIndex(0, runeLen)
			b := make([]byte, pwdLen)
			b[0] = runes[j]
			for i := 1; i < pwdLen; i++ {
				j = getIndex(j, runeLen)
				b[i] = runes[j]
			}
			sender <- string(b)
		}()
	}
	wgOuter.Wait()
}

func getIndex(last, l int) int {
	rand.Seed(seed())
	v := rand.Intn(l)
	if v == last {
		return getIndex(last, l)
	}
	return v
}

func seed() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return time.Now().UnixNano() + int64(m.Alloc) + int64(m.Frees)
}
