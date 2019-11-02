package main

import (
	"fmt"
	"github.com/jayvib/app/cmd/hacking/redis"
	"io/ioutil"
	"log"
	"os"
)

var client = redis.NewLocalClient()

func main() {
	setIncr()
}

func setIncr() {
	_, err := client.Set("counter", 100, 0).Result()
	handleError(err)

	_, err = client.Incr("counter").Result()
	handleError(err)
}

func setExample() {
	res, err := client.Set("hello", "jayson", 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func setFailIfExists() {
	res, err := client.SetNX("hello", "world", 0).Result()
	handleError(err)
	fmt.Println(res)
}

func setSuccessIfExists() {
	res, err := client.SetXX("hello", "world", 0).Result()
	handleError(err)
	fmt.Println(res)
}

func setImageValue() {
		file, err := os.Open("simple.png")
		handleError(err)
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		handleError(err)
		_, err = client.Set("simple.png", content, 0).Result()
		handleError(err)
}

func getImageValue() {
	content, err := client.Get("simple.png").Bytes()
	handleError(err)
	tmp, err := os.Create("simple-get.png")
	handleError(err)
	defer func() {
		err = tmp.Close()
		handleError(err)
	}()

	_, err = tmp.Write(content)
	handleError(err)
	fmt.Print("Success!")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
