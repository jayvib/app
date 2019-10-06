package main

import (
	"bytes"
	"testing"
)

// Requirements:
// Print 3
// Print 3 to Go!
// Wait a second between each line.

func TestCountdown(t *testing.T) {
	var buffer bytes.Buffer
	var spySleeper SpySleeper
	Countdown(&buffer, &spySleeper)
	got := buffer.String()
	want := `3
2
1
Go!`
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
	if spySleeper.Calls != 4 {
		t.Errorf("expecting '%d' calls but got '%d'", 4, spySleeper.Calls)
	}
}


type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

