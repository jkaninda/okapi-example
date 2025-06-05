package main

import (
	"github.com/jkaninda/okapi"
	"log/slog"
	"net/http"
)

var basicAuth = okapi.BasicAuthMiddleware{
	Username: "admin",
	Password: "password",
	Realm:    "Restricted Area",
}

func customMiddleware(next okapi.HandleFunc) okapi.HandleFunc {
	return func(c okapi.Context) error {
		slog.Info("Custom middleware executed", "path", c.Request.URL.Path, "method", c.Request.Method)
		// Call the next handler in the chain
		if err := next(c); err != nil {
			// If an error occurs, log it and return a generic error response
			slog.Error("Error in custom middleware", "error", err)
			return c.JSON(http.StatusInternalServerError, okapi.M{"error": "Internal Server Error"})
		}
		return nil
	}
}
