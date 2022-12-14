package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

var lock = sync.Mutex{}

func countLetters(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	for i := 0; i <= 20; i++ {
		for _, b := range body {
			c := strings.ToLower(string(b))
			index := strings.Index(allLetters, c)
			if index >= 0 {
				atomic.AddInt32(&frequency[index], 1)
			}
		}
	}
	wg.Done()
}

func main() {
	var frequency [26]int32
	wg := sync.WaitGroup{}
	start := time.Now()
	for i := 1000; i <= 1200; i++ {
		wg.Add(1)
		countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Printf("Processing took %s\n", elapsed)
	for i, f := range frequency {
		fmt.Printf("%s -> %d\n", string(allLetters[i]), f)
	}
}
