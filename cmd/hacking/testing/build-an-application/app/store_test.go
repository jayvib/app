package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	initialData := `[
	{"Name": "Luffy", "Wins": 10},
	{"Name": "Sanji", "Wins": 20}]`
	database, teardown := createTempFile(t, initialData)
	defer teardown()

	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)
	t.Run("league from a reader", func(t *testing.T) {
		want := League{
			{Name: "Sanji", Wins: 20},
			{Name: "Luffy", Wins: 10},
		}
		assertStoreGetLeague(t, store, want)

		t.Run("another read", func(t *testing.T) {
			assertStoreGetLeague(t, store, want)
		})
	})

	t.Run("get player score", func(t *testing.T) {
		assertStoreGetPlayerScore(t, store, "Luffy", 10)
	})

	t.Run("record wins", func(t *testing.T) {
		want := 11
		name := "Luffy"
		store.RecordWin(name)
		assertStoreGetPlayerScore(t, store , name, want)
		t.Run("record second time", func(t *testing.T){
			want = 12
			store.RecordWin(name)
			assertStoreGetPlayerScore(t, store, name, want)
		})
	})

	t.Run("store wins for the new players", func(t *testing.T){
		database, teardown := createTempFile(t, `[
		{"Name": "Guko", "Wins": 10},
		{"Name": "Vegita", "Wins": 20}]`)
		defer teardown()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		name := "Luffy"
		store.RecordWin(name)
		assertStoreGetPlayerScore(t, store, name, 1)
	})

	t.Run("works with an empty file", func(t *testing.T){
		db, teardown := createTempFile(t, "")
		defer teardown()

		_, err = NewFileSystemPlayerStore(db)
		assertNoError(t, err)
	})

	t.Run("sorted league", func(t *testing.T){
		tempFile, teardown := createTempFile(t, `[
{"Name": "Luffy", "Wins": 10},
{"Name": "Sanji", "Wins": 20}]`)
		defer teardown()

		store, err := NewFileSystemPlayerStore(tempFile)
		assertNoError(t, err)

		want := League{
			{Name: "Sanji", Wins: 20},
			{Name: "Luffy", Wins: 10},
		}

		assertStoreGetLeague(t, store, want)
		assertStoreGetLeague(t, store, want)
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expecting no error but got '%v'", err)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()
	tempFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}

	tempFile.WriteString(initialData)

	return tempFile, func() {
		err = tempFile.Close()
		if err != nil {
			t.Errorf("unable to close file: %v", err)
		}
		err := os.Remove(tempFile.Name())
		if err != nil {
			t.Errorf("unable to remove file: %v", err)
		}
	}
}

func assertStoreGetLeague(t *testing.T, store *FileSystemPlayerStore, want League) {
	t.Helper()
	got := store.GetLeague()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want '%v' got '%v'", want, got)
	}
}

func assertStoreGetPlayerScore(t *testing.T, store *FileSystemPlayerStore, name string, want int) {
	got := store.GetPlayerScore(name)
	if want != got {
		t.Errorf("want '%v' got '%v'", want, got)
	}
}
