package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Luffy")
	want := "Hello, Luffy"
	got := buffer.String()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
