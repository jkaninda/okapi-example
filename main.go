package main

import (
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/routes"
)

func main() {
	// Create a new Okapi instance with default config
	app := okapi.New()
	route := routes.NewRoute(app)

	// ************ Registering Routes ************
	// Register home route
	app.Register(route.Home())
	app.Register(route.WhoAmI())
	// Auth
	app.Register(route.AuthRoute())
	// Register book routes
	app.Register(route.BookRoutes()...)
	app.Register(route.CommonRoutes()...)
	// Admin routes
	app.Register(route.AdminRoutes()...)

	// Start the server
	if err := app.Start(); err != nil {
		panic(err)
	}
}
