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

package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/middlewares"
	"github.com/jkaninda/okapi-example/models"
	"net/http"
	"os"
	"strconv"
)

type BookController struct{}
type HomeController struct{}
type AuthController struct{}

var (
	books = []*models.Book{
		{Id: 1, Title: "Book One", Price: 100},
		{Id: 2, Title: "Book Two", Price: 200},
		{Id: 3, Title: "Book Three", Price: 300},
	}
)

// ****************** Controllers *****************

func (hc *HomeController) Home(c okapi.Context) error {
	return c.OK(okapi.M{"message": "Welcome to the Okapi Web Framework!"})
}

func (hc *HomeController) WhoAmI(c okapi.Context) error {
	email := c.Header("current_user_email")
	if email == "" {
		logger.Warn("no email found")
	}
	return c.OK(models.WhoAmIResponse{
		Host:   c.Request().Host,
		RealIp: c.RealIP(),
		CurrentUser: models.UserInfo{
			Name:  c.Header("current_user_name"),
			Email: email,
			Role:  c.Header("current_user_role"),
		},
	})
}
func (bc *BookController) GetBooks(c okapi.Context) error {
	_, err := bc.readBooksFromFile()
	if err != nil {
		logger.Error("Error reading books from file", "error", err)
		return c.ErrorInternalServerError(models.ErrorResponse{Success: false, Status: http.StatusInternalServerError, Details: err.Error()})
	}
	return c.OK(books)
}

func (bc *BookController) CreateBook(c okapi.Context) error {
	// Simulate creating a book in a database
	book := &models.Book{}
	err := c.Bind(book)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse{Success: false, Status: http.StatusBadRequest, Details: err.Error()})
	}
	book.Id = len(books) + 1
	books = append(books, book)
	response := models.Response{
		Success: true,
		Message: "Book created successfully",
		Data:    *book,
	}
	return c.OK(response)
}
func (bc *BookController) GetBook(c okapi.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse{Success: false, Status: http.StatusBadRequest, Details: err.Error()})
	}
	// Simulate a fetching book from a database

	_, err = bc.readBooksFromFile()
	if err != nil {
		logger.Error("Error reading books from file", "error", err)
		return c.ErrorInternalServerError(models.ErrorResponse{Success: false, Status: http.StatusInternalServerError, Details: err.Error()})
	}
	for _, book := range books {
		if book.Id == i {
			return c.OK(book)
		}
	}
	return c.AbortNotFound("Book not found")
}

// ******************** AuthController *****************

func (bc *AuthController) Login(c okapi.Context) error {
	authRequest := &models.AuthRequest{}
	err := c.Bind(authRequest)
	if err != nil {
		return c.ErrorBadRequest(models.ErrorResponse{Success: false, Status: http.StatusBadRequest, Details: err.Error()})
	}
	// Validate the authRequest and generate a JWT token
	authResponse, err := middlewares.Login(authRequest)
	if err != nil {
		return c.ErrorUnauthorized(authResponse)
	}
	return c.OK(authResponse)
}
func (bc *AuthController) WhoAmI(c okapi.Context) error {
	//Get User Information from the context, shared by the JWT middleware using forwardClaims
	email := c.GetString("email")
	if email == "" {
		return c.AbortUnauthorized("Unauthorized", fmt.Errorf("user not authenticated"))
	}

	c.Response().Header().Set("X-Okapi-User", email)
	c.Response().Header().Set("X-Okapi-User-Name", c.GetString("name"))
	c.Response().Header().Set("X-Okapi-Role", c.GetString("role"))
	// Respond with the current user information
	return c.OK(models.UserInfo{
		Email: email,
		Role:  c.GetString("role"),
		Name:  c.GetString("name"),
	},
	)
}

func (bc *BookController) readBooksFromFile() ([]*models.Book, error) {
	if books != nil && len(books) > 3 {
		return books, nil
	}
	booksFile, err := os.ReadFile("data/books.json")
	if err != nil {
		logger.Error("Error reading books file", "error", err)
		return nil, fmt.Errorf("failed to read books data: %w", err)
	}
	err = json.Unmarshal(booksFile, &books)
	if err != nil {
		logger.Error("Error unmarshalling books data", "error", err)
		return nil, fmt.Errorf("failed to parse books data: %w", err)
	}
	return books, nil
}
