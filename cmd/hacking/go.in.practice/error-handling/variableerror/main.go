package main

import (
	"errors"
	"fmt"
)

var ErrLuffy = errors.New("luffy can't mate with the girl")
var ErrSanji = errors.New("sanji can't mate with the girl")

func mateMatching(girlName, boyName string) (string, error) {
	if boyName == "luffy" {
		return "", ErrLuffy
	} else if boyName == "sanji" {
		return "", ErrSanji
	} else {
		return "okeee keeyoohhh", nil
	}
}

func main() {
	fmt.Println(mateMatching("althea", "luffy"))
	fmt.Println(mateMatching("althea", "sanji"))
	fmt.Println(mateMatching("althea", "jayson"))
}
