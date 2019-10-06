package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	countdownStart = 3
	finalWord = "Go!"
)

type DefaultSleeper struct {}

func (DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

type Sleeper interface {
	Sleep()
}

func Countdown(w io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		fmt.Fprintf(w, "%d\n", i)
	}
	sleeper.Sleep()
	fmt.Fprintf(w, finalWord)
}

func main() {
	Countdown(os.Stdout, DefaultSleeper{})
}
