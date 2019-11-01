package book_test

import (
	"testing"

	. "github.com/jayvib/app/cmd/hacking/ginkgo/book"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book", func() {
	var (
		longBook  Book
		shortBook Book
	)

	// Set up state for the specs.
	BeforeEach(func() {
		longBook = Book{
			Title:  "One Piece",
			Author: "Luffy MOnkey",
			Pages:  1222,
		}

		shortBook = Book{
			Title:  "Fairy Tail",
			Author: "Natsu Dragneel",
			Pages:  280,
		}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})

		Context("With fewer than 300 pages", func() {
			// It a descriptive assertion
			It("should be a short story", func() {
				// Expect is for expectation
				// Equal is the assertion
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
	})

	Describe("Load from JSON", func() {
		It("can be loaded from JSON", func() {
			book, _ := NewBookFromJSON(`{
			"title": "One Piece",
			"author": "Luffy Monkey",
			"pages": 1222
		}`)
			Expect(book.Title).To(Equal("One Piece"))
			Expect(book.Author).To(Equal("Luffy Monkey"))
			Expect(book.Pages).To(Equal(1222))
		})

		Context("when calling NewBookFromJSON()", func() {
			Context("when empty string provided", func() {
				Specify("an ErrEmptyString error is returned", func() {
					_, err := NewBookFromJSON("")
					Expect(err).To(Equal(ErrEmptyString))
				})
			})
		})
	})

})

func TestBook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Book Suite")
}
