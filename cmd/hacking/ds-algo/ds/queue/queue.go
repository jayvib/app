package main

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Person struct {
	Name string
}

func main() {
	q := generateInt(1000)

	var wg sync.WaitGroup
	workers := 10
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sink(q, id+1)
		}(i)
	}

	wg.Wait()
}

func sink(q *queue.Queue, id int) {
	batch := 0
	for {
		batch++
		val, err := q.Poll(5, time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d: %v\n", id,val)
	}
}

func generateInt(n int) (q *queue.Queue) {
	q = queue.New(int64(n))
	go func() {
		for i := 0; i < n; i++ {
			time.Sleep(100 * time.Millisecond)
			num := rand.Int63()
			err := q.Put(num)
			handleErr(err)
		}
	}()
	return
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func process(d interface{}) (data interface{}) {
	switch v := d.(type) {
	case int:
		return v * 2
	case string:
		return strings.ToUpper(v)
	}
	return d
}
