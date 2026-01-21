package services

import (
	"strconv"
	"strings"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/data"
	"github.com/jkaninda/okapi-example/models"
)

type BookService struct{}

var (
	books = []*models.Book{}
)

func (bc *BookService) List(c *okapi.Context) error {
	q := c.Query("q")
	books = bc.booksData()
	if len(q) > 0 {
		return c.OK(searchBooks(books, q))
	}
	return c.OK(books)
}

func (bc *BookService) Create(c *okapi.Context) error {
	// Simulate creating a book in a database
	book := &models.Book{}
	err := c.Bind(book)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse)
	}
	book.Id = len(books) + 1
	books = append(books, book)

	return c.OK(models.SuccessResponse("Book created successfully", book))
}
func (bc *BookService) Get(c *okapi.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse("Bad Request", err))
	}
	// Simulate a fetching book from a database

	books = bc.booksData()

	for _, book := range books {
		if book.Id == i {
			return c.OK(book)
		}
	}
	return c.AbortNotFound("Book not found")
}
func (bc *BookService) booksData() []*models.Book {
	books, _ = data.Books()
	return books
}
func searchBooks(books []*models.Book, query string) []*models.Book {
	if query == "" {
		return books
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []*models.Book

	for _, book := range books {
		titleMatch := strings.Contains(strings.ToLower(book.Title), query)
		authorMatch := strings.Contains(strings.ToLower(book.Author), query)

		if titleMatch || authorMatch {
			results = append(results, book)
		}
	}

	return results
}
