package book

import (
	"encoding/json"
	"errors"
)

var ErrEmptyString = errors.New("empty string value")

type Book struct {
	Title  string
	Author string
	Pages  int
}

func (b *Book) CategoryByLength() string {

	if b.Pages > 300 {
		return "NOVEL"
	} else {
		return "SHORT STORY"
	}
}

func NewBookFromJSON(s string) (*Book, error) {
	if s == "" {
		return nil, ErrEmptyString
	}
	var book Book
	json.Unmarshal([]byte(s), &book)
	return &book, nil
}
