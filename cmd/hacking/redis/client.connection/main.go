package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v7"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)
}
