package main

import "fmt"

// WISDOM:
// - It is good to separate your "domain" code from the outside world
// - Refactoring! It is important that your tests are clear specifications
//   of what the code needs to do.

const (
	englishPrefix = "Hello, "
	spanishPrefix = "Hola, "
	frenchPrefix  = "Bonjour, "
)

const (
	defaultToGreet = "World"
)

const (
	spanish = "Spanish"
	french  = "French"
)

// Hello takes a name and returns a
// greeting message.
func Hello(name string, language string) string {
	if name == "" {
		name = defaultToGreet
	}
	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		return spanishPrefix
	case french:
		return frenchPrefix
	default:
		return englishPrefix
	}
}

func main() {
	fmt.Println(Hello("pocker", ""))
}
