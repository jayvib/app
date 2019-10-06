package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

// Requirements:
// Print 3
// Print 3 to Go!
// Wait a second between each line.

func TestCountdown(t *testing.T) {
	t.Run("countdown from 3 to Go!", func(t *testing.T) {
		var buffer bytes.Buffer
		spySleeper := &CoundownOperationSpy{}
		Countdown(&buffer, spySleeper)
		got := buffer.String()
		want := `3
2
1
Go!`
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("sleep after every print", func(t *testing.T){
		spySleepWriter := &CoundownOperationSpy{}
		Countdown(spySleepWriter, spySleepWriter)
		want := []string{
			sleep, write,
			sleep, write,
			sleep, write,
			sleep, write,
		}
		if !reflect.DeepEqual(spySleepWriter.Calls, want) {
			t.Errorf("wanted calls %v got %v", want, spySleepWriter.Calls)
		}
	})
}

func TestConfigurableSleeper(t *testing.T) {
	spyTime := &SpyTime{}
	sleepTime := 5 * time.Second
	configSleeper := &ConfigurableSleeper{duration: sleepTime, sleep: spyTime.Sleep}
	configSleeper.Sleep()
	got := spyTime.durationSlept
	want := sleepTime
	if got != want {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(d time.Duration) {
	s.durationSlept = d
}

const (
	sleep = "sleep"
	write = "write"
)

type CoundownOperationSpy struct {
	Calls []string
}

func (c *CoundownOperationSpy) Sleep() {
	c.Calls = append(c.Calls, sleep)
}

func (c *CoundownOperationSpy) Write(p []byte) (n int, err error){
	c.Calls = append(c.Calls, write)
	return len(p), nil
}


type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

