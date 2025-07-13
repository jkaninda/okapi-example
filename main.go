package main

import (
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/routes"
)

func main() {
	// Create a new Okapi instance with default config
	app := okapi.Default()
	// ************ Registering Routes ************
	// Register home route
	app.Register(routes.Home())
	app.Register(routes.WhoAmI())
	// Auth
	app.Register(routes.AuthRoute())
	// Register book routes
	app.Register(routes.BookRoutes()...)
	app.Register(routes.CommonRoutes()...)
	// Admin routes
	app.Register(routes.AdminRoutes()...)

	// Start the server
	if err := app.Start(); err != nil {
		panic(err)
	}
}
