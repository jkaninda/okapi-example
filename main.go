/*
 *  MIT License
 *
 * Copyright (c) 2025 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"net/http"
)

type (
	User struct {
		ID       int    `param:"id" query:"id" form:"id" json:"id" xml:"id" max:"10" `
		Name     string `json:"name" form:"name"  max:"15"`
		IsActive bool   `json:"is_active" query:"is_active" yaml:"isActive"`
	}
	Book struct {
		ID    int    `json:"id" param:"id" query:"id" form:"id" xml:"id" max:"50" `
		Name  string `json:"name" form:"name"  max:"50" default:"anonymous"`
		Price int    `json:"price" form:"price" query:"price" yaml:"price" `
		Qty   int    `json:"qty" form:"qty" query:"qty" yaml:"qty"`
	}
)

var (
	books = []*Book{
		{ID: 1, Name: "The Go Programming Language", Price: 30, Qty: 100},
		{ID: 2, Name: "Learning Go", Price: 25, Qty: 50},
		{ID: 3, Name: "Go in Action", Price: 40, Qty: 75},
		{ID: 4, Name: "Go Web Programming", Price: 35, Qty: 60},
		{ID: 5, Name: "Go Design Patterns", Price: 45, Qty: 80},
	}
)

var (
	users = []*User{
		{ID: 1, Name: "John Doe", IsActive: true},
		{ID: 2, Name: "Jonas Kaninda", IsActive: false},
		{ID: 3, Name: "Alice Johnson", IsActive: true},
	}
)

func main() {
	log := logger.New(logger.WithCaller(), logger.WithJSONFormat())
	// Create a new Okapi instance and use custom logger
	o := okapi.New(okapi.WithReadTimeout(15), okapi.WithWriteTimeout(15), okapi.WithLogger(log.Logger))
	// Enable openapi documentation
	o.With().WithOpenAPIDocs()

	o.Get("/", func(c okapi.Context) error {
		// Handler logic for the root route
		return c.OK(okapi.M{"message": "Welcome to Okapi!"})
	}, okapi.DocSummary("Welcome page"), okapi.DocResponse(okapi.M{}))
	// Create a new group with a base path
	api := o.Group("/api")
	// Create and apply custom middleware to the v1 group
	v1 := api.Group("/v1", customMiddleware)

	// Define a route with a handler
	v1.Get("/users", func(c okapi.Context) error {
		return c.OK(users)
	},
		okapi.Doc().Summary("Get all users").Response([]User{}).Build(),
	)
	// Get user
	v1.Get("/users/:id", show, okapi.DocPathParam("id", "int", "User Id"))
	// Update user
	v1.Put("/users/:id", update, okapi.DocPathParam("id", "int", "User Id"))
	// Create user
	v1.Post("/users", store, okapi.DocRequestBody(User{}), okapi.DocResponse(User{}))

	// Create a new group with a base path v2
	v2 := api.Group("/v2")
	// Define a route with a handler
	v2.Get("/users", func(c okapi.Context) error {
		c.SetHeader("Version", "v2")
		// Handler logic for the route
		return c.OK(users)
	},
		okapi.Doc().Summary("Get all users").Response([]User{}).Build(),
	)
	v2.Get("/users/:id", func(c okapi.Context) error {
		c.SetHeader("Version", "v2")
		id := c.Param("id")
		for _, user := range users {
			if fmt.Sprintf("%d", user.ID) == id {
				return c.JSON(http.StatusOK, user)
			}
		}
		return c.ErrorNotFound(map[string]string{"error": "User not found"})
	},
		okapi.Doc().Summary("Get user by Id").Response([]User{}).
			PathParam("id", "int", "User Id").
			Build(),
	)

	// Create a new group with a base path for admin routes and apply basic auth middleware
	adminApi := api.Group("/admin", basicAuth.Middleware) // This group will require basic authentication
	adminApi.Put("/books/:id", adminUpdate, okapi.DocPathParam("id", "int", "Book Id"), okapi.DocResponse(Book{}))
	adminApi.Post("/books", adminStore, okapi.DocRequestBody(Book{}), okapi.DocResponse(Book{}))

	// Define routes for the v1 group
	v1.Get("/books", getBooks, okapi.DocResponse([]Book{}))
	v1.Get("/books/:id", getBook, okapi.DocPathParam("id", "int", "Book Id"), okapi.DocResponse(Book{})).Name = "show_book" // Named route for easier reference

	// Start the server
	if err := o.Start(); err != nil {
		panic(err)
	}

}

// ***** Handlers *****
func store(c okapi.Context) error {
	var newUser User
	if ok, err := c.ShouldBind(&newUser); !ok {
		errMessage := fmt.Sprintf("Failed to bind user data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input " + errMessage})
	}

	// Add the new user to the list
	newUser.ID = len(users) + 1 // Simple ID assignment
	users = append(users, &newUser)
	// Respond with the created user
	return c.JSON(http.StatusCreated, newUser)
}
func show(c okapi.Context) error {
	var newUser User
	if ok, err := c.ShouldBind(&newUser); !ok {
		errMessage := fmt.Sprintf("Failed to bind user data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input " + errMessage})
	}
	for _, user := range users {
		if user.ID == newUser.ID {
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNotFound, okapi.M{"error": "User not found"})
}
func update(c okapi.Context) error {
	var newUser User
	if ok, err := c.ShouldBind(&newUser); !ok {
		errMessage := fmt.Sprintf("Failed to bind user data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input " + errMessage})
	}
	for _, user := range users {
		if user.ID == newUser.ID {
			user.Name = newUser.Name
			user.IsActive = newUser.IsActive
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNotFound, okapi.M{"error": "User not found"})
}

func adminStore(c okapi.Context) error {
	var newBook Book
	if ok, err := c.ShouldBind(&newBook); !ok {
		errMessage := fmt.Sprintf("Failed to bind book: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input " + errMessage})
	}
	// Get username
	username := c.GetString("username")
	fmt.Printf("Current user: %s\n", username)
	// Add the new book to the list
	newBook.ID = len(books) + 1 // Simple ID assignment
	books = append(books, &newBook)
	// Respond with the created book
	return c.JSON(http.StatusCreated, newBook)
}
func adminUpdate(c okapi.Context) error {
	var newBook Book
	if ok, err := c.ShouldBind(&newBook); !ok {
		return c.ErrorBadRequest(err)
	}
	for _, book := range books {
		if book.ID == newBook.ID {
			book.Name = newBook.Name
			book.Price = newBook.Price
			book.Qty = newBook.Qty
			// Respond with the updated book
			return c.JSON(http.StatusOK, book)
		}
	}
	return c.ErrorNotFound("Book not found")
}
func getBooks(c okapi.Context) error {
	return c.OK(books)
}
func getBook(c okapi.Context) error {
	var newBook Book
	// Bind the book ID from the request parameters using `param` tags
	// You can also use c.Param("id") to get the ID from the URL
	if ok, err := c.ShouldBind(&newBook); !ok {
		errMessage := fmt.Sprintf("Failed to bind book: %v", err)
		return c.ErrorBadRequest(map[string]string{"error": "Invalid input " + errMessage})
	}
	// time.Sleep(2 * time.Second) // Simulate a delay for demonstration purposes

	for _, book := range books {
		if book.ID == newBook.ID {
			return c.OK(book)
		}
	}
	return c.ErrorNotFound(okapi.M{"error": "Book not found"})
}
