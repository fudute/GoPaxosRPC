package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func randString(n int) string {
	lenght := rand.Int()%n + 1

	str := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		str[i] = byte(rand.Int()%26) + 'a'
	}
	return string(str)
}

func testGet(n int) {
	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go func() {
			_, err := http.Get("http://127.0.0.1:8000/get/name")
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println(time.Since(start))

}

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func testSet(n int) {

	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {

		kvp := &KVPair{
			Key:   randString(10),
			Value: randString(10),
		}

		bs, err := json.Marshal(kvp)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * 1)
		go func() {
			_, err := http.Post("http://127.0.0.1:8000/store/set", "application/json", bytes.NewReader(bs))
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println(time.Since(start))
}

// "set key value"
func main() {
	n := 100
	testSet(n)
}
