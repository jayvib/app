package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func main() {
	ReadRuneExample()
}

func ReadRuneExample() {
	reader := bufio.NewReader(strings.NewReader("hello"))

	ReadRuneAndPrint(reader)
	fmt.Println("Unreading last rune")
	err := reader.UnreadRune()
	if err != nil {
		log.Fatal(err)
	}
	ReadRuneAndPrint(reader)
}

func ReadRuneAndPrint(reader *bufio.Reader) {
	c, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Read rune:", string(c))
}
