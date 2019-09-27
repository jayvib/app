package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFun1(t *testing.T) {
	input := 10
	got := Fun1(input)
	want := input
	assert.Equal(t, got, want)
}

func TestFun2(t *testing.T) {
	input := 10
	got := Fun2(input)
	want := pow(input, 2)
	assert.Equal(t, got, want)
}

func TestFun4(t *testing.T) {
	input := 10
	got := Fun4(input)
	want := pow(input, 3)
	assert.Equal(t, want, got)
}

func pow(n,  p int) int {
	return int(math.Pow(float64(n), float64(p)))
}