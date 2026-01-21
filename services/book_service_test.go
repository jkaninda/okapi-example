package services

import (
	"testing"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi/okapitest"
)

var bookService = &BookService{}

func TestBooksAPI(t *testing.T) {
	// Setup test server
	server := okapi.NewTestServer(t)

	server.Get("/v1/books", bookService.List)
	server.Get("/v1/books/{id:int}", bookService.Get)

	// Create reusable client
	client := okapitest.NewClient(t, server.BaseURL)

	client.GET("/v1/books").ExpectStatusOK().ExpectBodyContains("The Kubernetes Bible")
	client.GET("/v1/books?q=Vol2").ExpectStatusOK().ExpectBodyContains("System Design Interview Vol2")
	client.GET("/v1/books/4").ExpectStatusOK().ExpectBodyContains("System Design Interview Vol2")

}
