package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	countdownStart = 3
	finalWord      = "Go!"
)

type Sleeper interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(duration time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
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
	configSleeper := &ConfigurableSleeper{
		duration: 1 * time.Second,
		sleep: time.Sleep,
	}
	Countdown(os.Stdout, configSleeper)
}
