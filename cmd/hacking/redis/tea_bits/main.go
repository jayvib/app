package main

import (
	"fmt"
	"github.com/jayvib/app/cmd/hacking/redis"
)

func populateTea() {
	client := redis.NewLocalClient()

	for i := 0; i < 6000; i++ {
		key := fmt.Sprintf("tea/%d", i)
		client.SAdd("teas/caffeinated", key)
		client.SetBit("teas/caffeine", int64(i), 1)
	}

	err := client.FlushAll()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	populateTea()
}
