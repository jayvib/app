package walk

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {

	tests := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"Struct with one string field",
			struct{ Name string }{"Luffy"},
			[]string{"Luffy"},
		},
		{
			"Struct with one string field",
			struct {
				Name string
				City string
			}{"Luffy", "East Blue"},
			[]string{"Luffy", "East Blue"},
		},
		{
			"Struct with non-string field",
			struct {
				Name string
				Age  int
			}{"Luffy", 20},
			[]string{"Luffy"},
		},
		{
			"Nested struct",
			Person{
				Name: "Luffy",
				Profile: Profile{ Age: 22, City: "EastBlue" },
			},
			[]string{"Luffy", "EastBlue"},
		},
		{
			"Pointers to things",
			&Person{
				"Luffy",
				Profile{22, "EastBlue"},
			},
			[]string{"Luffy", "EastBlue"},
		},
		{
			"Slices",
			[]Profile{
				{22, "EastBlue"},
				{23, "WestBlue"},
			},
			[]string{"EastBlue", "WestBlue"},
		},
		{
			"Arrays",
			[2]Profile{
				{22, "EastBlue"},
				{23, "WestBlue"},
			},
			[]string{"EastBlue", "WestBlue"},
		},
		{
			"Maps",
			map[string]string{
				"Foo": "Bar",
				"Baz": "Boz",
			},
			[]string{"Bar", "Boz"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			var got []string
			walk(tc.Input, func(s string) {
				got = append(got, s)
			})

			if !reflect.DeepEqual(tc.ExpectedCalls, got) {
				t.Errorf("wrong number of calls: want '%v' got '%v'", tc.ExpectedCalls, got)
			}
		})
	}

}
