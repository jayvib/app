package syncpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	p := &Person{
		Name: "Luffy Monkey",
	}

	got := MarshalJSON(p)
	want := `{"name":"Luffy Monkey"}
`

	assert.Equal(t, want, got)
}

var result string

// Without sync pool:
// BenchmarkMarshalJSON-4   	 2188296	       494 ns/op
func BenchmarkMarshalJSON(b *testing.B) {
	p := &Person{
		Name: "Luffy Monkey",
	}
	var r string
	for i := 0; i < b.N; i++ {
		// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
		r = MarshalJSON(p)
	}
	result = r
}