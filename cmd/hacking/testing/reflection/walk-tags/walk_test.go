package walk_tags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name string `whome:"name"`
	Profile Profile
}

type Profile struct {
	Age int `whome:"age"`
	City string `whome:"city"`
	Job string `whome:"job"`
}

// Walk through all the fields that has a tags name "whome"
// and store the value into a map with its equivalent tag value
// and the field value.
func TestWalk(t *testing.T) {

	tests := []struct{
		name string
		input interface{}
		expected map[string]interface{}
	}{
		{
			name: "single field struct",
			input: struct{ Name string `whome:"name"`}{ Name: "Luffy" },
			expected: map[string]interface{}{ "name": "Luffy" },
		},
		{
			name: "two field struct",
			input: struct{
				Name string `whome:"name"`
				Age int `whome:"age"`
			}{
				Name: "Luffy",
				Age: 22,
			},
			expected: map[string]interface{}{ "name": "Luffy", "age": 22},
		},
		{
			name: "nested struct",
			input: Person{
				Name: "Luffy",
				Profile: Profile{
					Age: 22,
					City: "EastBlue",
					Job: "Pirate",
				},
			},
			expected: map[string]interface{}{
				"name": "Luffy",
				"age": 22,
				"city": "EastBlue",
				"job": "Pirate",
			},
		},
		{
			name: "pointer struct",
			input: &Person{
				Name: "Luffy",
				Profile: Profile{
					Age: 22,
					City: "EastBlue",
					Job: "Pirate",
				},
			},
			expected: map[string]interface{}{
				"name": "Luffy",
				"age": 22,
				"city": "EastBlue",
				"job": "Pirate",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T){
			got := make(map[string]interface{})
			walk(tc.input, func(key string, value interface{}){
				got[key] = value
			})

			assert.Equal(t, tc.expected, got)
		})
	}

}
