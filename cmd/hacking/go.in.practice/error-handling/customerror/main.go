package main

import "fmt"

const CodeDivZero = iota

type DivError struct {
	errCode int
}

func (d *DivError) Error() string {
	switch d.errCode {
	case CodeDivZero:
		return "cannot divide b zero digit"
	}
	return ""
}

func quotient(n, d int) (int, error) {
	if d == 0 {
		return 0, &DivError{errCode: CodeDivZero}
	}
	return n / d, nil
}

func main() {
	_, err := quotient(2, 0)
	if err != nil {
		if _, ok := err.(*DivError); ok {
			fmt.Println("oppps division errorr")
		}
	}

}
