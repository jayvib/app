package app

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	// test when the content to write is smaller
	// than the current content in the writer

	file, teardown := createTempFile(t, "12345")
	defer teardown()

	store := &tape{file}

	_, err := store.Write([]byte("abc"))
	if err != nil {
		t.Fatal(err)
	}

	want := "abc"

	file.Seek(0, io.SeekStart)
	content, _ := ioutil.ReadAll(file)
	got := string(content)
	if got != want {
		t.Errorf("want '%s' got '%s'", want, got)
	}
}
