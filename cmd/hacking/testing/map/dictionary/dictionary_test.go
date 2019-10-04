package dictionary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "this is just a test"}
	t.Run("known word", func(t *testing.T){
		key := "test"
		want := "this is just a test"
		assertDefinition(t, dictionary, key, want)
	})

	t.Run("unknown word", func(t *testing.T){
		_, err := dictionary.Search("unknown")
		assertError(t, err, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T){
		dictionary := Dictionary{}
		key := "test"
		want := "this is just a test"
		dictionary.Add(key, want)
		assertDefinition(t, dictionary, key, want)
	})

	t.Run("existing word", func(t *testing.T){
		key := "test"
		want := "this is just a test"
		dictionary := Dictionary{key: want}
		err := dictionary.Add(key, want)
		assertDefinition(t, dictionary, key, want)
		if assert.IsType(t, Err{}, err) {
			e := err.(Err)
			assert.Equal(t, WordExist, e.t)
			assert.Equal(t, Op("Dictionary.Add"), e.operation)
			t.Log(e.Error())
		}
	})

}

func assertDefinition(t *testing.T, d Dictionary, key, want string) {
	t.Helper()
	got, err := d.Search(key)
	assertNoError(t, err)
	assertString(t, got, want)
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expecting no error but got one '%v'", err)
	}
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("expecting error '%v' got '%v'", want, got)
	}
}

