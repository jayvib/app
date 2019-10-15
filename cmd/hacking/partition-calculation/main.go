package main

import (
	"crypto/md5"
	"fmt"
	"math/big"
)

func main() {
	key := "abc"
	partition := int64(2)
	fmt.Println(Partition(key, partition))
	key = "abcd"
	fmt.Println(Partition(key, partition))
	key = "abcdefg"
	fmt.Println(Partition(key, partition))
	key = "abc232defg"
	fmt.Println(Partition(key, partition))
	key = "abc2323232defg"
	fmt.Println(Partition(key, partition))
}

func Partition(key string, count int64) int64 {
	hash := md5.New()
	hash.Write([]byte(key))
	i := new(big.Int)
	fmt.Println("i:", i)
	fmt.Println("hash sum:", string(hash.Sum(nil)))
	i = i.SetBytes(hash.Sum(nil))
	fmt.Println(i)
	fmt.Println(i.Int64() % count)
	return i.Mod(i, big.NewInt(count)).Int64()
}
