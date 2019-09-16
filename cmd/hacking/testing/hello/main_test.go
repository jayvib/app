package main

import "testing"

// TDD Cycle:
// 1. Write a test.
// 2. Make the compiler pass
// 3. Run the test, see that it fails and check the error message is meaningful
// 4. Write minimal code to pass the test.
// 5. Refactor.

// Quote:
// "By not writing tests you are committing to manually checking your code by
//	running your software which breaks your state of flow and you won’t be saving
//	yourself any time, especially in the long run."

// Case 1: Greet the person.
// Case 2: When a function is called with an empty string it defaults to "Hello, World"
// Case 3: If a language is passed in that we do not recognise, just default to English.
func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Luffy", "")
		want := "Hello, Luffy"
		assertString(t, got, want)
	})

	t.Run("When passed with an empty string it defaults to 'Hello, World", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertString(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Luffy", "Spanish")
		want := "Hola, Luffy"
		assertString(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Luffy", "French")
		want := "Bonjour, Luffy"
		assertString(t, got, want)
	})
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
