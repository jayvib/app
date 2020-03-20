package integers

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdder(t *testing.T) {
	t.Run("add even ", func(t *testing.T) {
		got, err := Add(2, 2)
		assert.NoError(t, err)
		assertSum(t, got, 4)
	})

	t.Run("adding number that is odd", func(t *testing.T){
		_, err := Add(1, 2)
		if err == nil {
			t.Fatal("expecting an error but got nil")
		}

		berr, ok := err.(BadParameter)
		if !ok {
			t.Fatalf("expecting a BadParameter type error but didn't got one")
		}

		if berr.Type != OddParameter {
			t.Errorf("want bad parameter type '%s' but got '%s'", OddParameter, berr.Type)
		}

		err = errors.Unwrap(err)
		if err != ErrParameterOdd {
			t.Errorf("want ErrParameterOdd but got '%s'", err)
		}

	})
}

func ExampleAdd() {
	sum, _ := Add(2, 2)
	fmt.Println(sum)
	// Output: 4
}

func assertSum(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got '%d' want '%d'\n", got, want)
	}
}