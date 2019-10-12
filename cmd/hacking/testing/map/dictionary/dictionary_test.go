package dictionary

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := &Dictionary{
		dictionary: map[string]string{"test": "this is just a test"},
	}
	t.Run("known word", func(t *testing.T) {
		key := "test"
		want := "this is just a test"
		assertDefinition(t, dictionary, key, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dictionary.Search("unknown")
		if assert.IsType(t, Err{}, got) {
			e := got.(Err)
			assert.Equal(t, WordNotExist, e.t)
			assert.Equal(t, Op("dictionary/Dictionary.Search"), e.operation)
			want := Err{t: WordNotExist, operation: Op("dictionary/Dictionary.Search")}
			assertError(t, want, got)
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dictionary := &Dictionary{
			dictionary: make(map[string]string),
		}
		key := "test"
		want := "this is just a test"
		dictionary.Add(key, want)
		assertDefinition(t, dictionary, key, want)
	})

	t.Run("existing word", func(t *testing.T) {
		key := "test"
		want := "this is just a test"
		dictionary := &Dictionary{
			dictionary: map[string]string{key: want},
		}
		got := dictionary.Add(key, want)
		assertDefinition(t, dictionary, key, want)
		if assert.IsType(t, Err{}, got) {
			wantErr := Err{WordExist, Op("dictionary/Dictionary.Add"), ErrWordExist}
			assertError(t, wantErr, got)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("word exists", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := &Dictionary{
			dictionary: map[string]string{
				word: definition,
			},
		}
		newDefinition := "new definition"
		dictionary.Update(word, newDefinition)
		assertDefinition(t, dictionary, word, newDefinition)
	})

	t.Run("word not exists", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{}
		got := dictionary.Update(word, definition)
		if assert.IsType(t, Err{}, got) {
			wantErr := Err{WordNotExist, Op("dictionary/Dictionary.Update"), ErrNotFound}
			assertError(t, wantErr, got)
		}
	})

	t.Run("multiple updates to the word", func(t *testing.T) {
		word := "test"
		definition := "this is a test for concurrency"

		dictionary := &Dictionary{
			dictionary: map[string]string{
				word: definition,
			},
		}

		numUpdater := 5
		var wg sync.WaitGroup

		newDefinition := "this is a new definition for test in concurrency"

		for i := 0; i < numUpdater; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				dictionary.Update(word, newDefinition)
			}(&wg)
		}

		wg.Wait()
		assertDefinition(t, dictionary, word, newDefinition)
	})
}

func TestDelete(t *testing.T) {
	word := "test"
	dictionary := &Dictionary{
		dictionary: map[string]string{
			word: "test definition",
		},
	}
	dictionary.Delete(word)

	_, err := dictionary.Search(word)
	if assert.IsType(t, Err{}, err) {
		wantErr := Err{t: WordNotExist, operation: Op("dictionary/Dictionary.Search")}
		assertError(t, wantErr, err)
	}
}

func assertDefinition(t *testing.T, d *Dictionary, key, want string) {
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
	switch g := got.(type) {
	case Err:
		w, ok := want.(Err)
		if !ok {
			t.Errorf("expecting error '%v' got '%v'", w, g)
		}
		assert.Equal(t, w.operation, g.operation)
		assert.Equal(t, w.t, g.t)
	default:
		if got != want {
			t.Errorf("expecting error '%v' got '%v'", want, got)
		}
	}
}
