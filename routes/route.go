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

package routes

import (
	"net/http"

	"github.com/jkaninda/okapi-example/config"
	"github.com/jkaninda/okapi-example/models"
	"github.com/jkaninda/okapi-example/services"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/middlewares"
)

var (
	bookService        = &services.BookService{}
	commonService      = &services.CommonService{}
	chatRoomService    = &services.ChatRoomService{}
	authService        = &services.AuthService{}
	bearerAuthSecurity = []map[string][]string{
		{
			"bearerAuth": {},
		},
	}
)

// You can also use this example

type Router struct {
	app    *okapi.Okapi
	group  *okapi.Group
	cfg    *config.Config
	assets http.FileSystem
}

// NewRoute creates a new Route instance with the provided Okapi app
func New(app *okapi.Okapi, conf *config.Config, assets http.FileSystem) *Router {
	commonService.SessionManager = conf.SessionManager
	return &Router{
		app:    app,
		group:  &okapi.Group{Prefix: "api/v1"},
		cfg:    conf,
		assets: assets,
	}
}

func (r *Router) RegisterRoutes() {
	r.registerAll()
	r.app.Register(r.whoAmI())
	r.app.Register(r.authRoute())
	r.app.Register(r.coreRoutes()...)
	r.app.Register(r.bookRoutes()...)
	r.app.Register(r.v1BookRoutes()...)
	// Admin routes
	r.app.Register(r.AdminRoutes()...)
}

// Home return Render
func (r *Router) registerAll() {
	r.app.Get("/", func(c *okapi.Context) error {
		return commonService.Home(c)
	})
	r.app.Get("/sse/sessions", func(c *okapi.Context) error {
		return commonService.Session(c)
	},
		okapi.Summary("Get current server time"),
	)
	r.app.Get("/chat", chatRoomService.ChatPage)

	// WebSocket endpoint
	r.app.Get("/ws", chatRoomService.WebSocketHandle,
		okapi.Summary("Start Websocket"),
		okapi.DocQueryParam("token", "string", "Websocket auth token", false),
	)
	// Static
	r.app.StaticFS("/assets", r.assets)

}

// ****************** Route Definitions ******************

// WhoAmI returns the route definition for the HomeController
func (r *Router) whoAmI() okapi.RouteDefinition {
	return okapi.RouteDefinition{
		Path:        "/whoami",
		Method:      http.MethodGet,
		Handler:     commonService.WhoAmI,
		Group:       r.group,
		Summary:     "Whoami",
		Description: "Get the current user's information, no auth requir",
		Request:     &models.WhoAmIRequest{},
		Response:    &models.WhoAmIResponse{},
	}
}

// ************* Book Routes *************

// bookRoutes returns the route definitions for the BookService
func (r *Router) bookRoutes() []okapi.RouteDefinition {
	apiGroup := &okapi.Group{Prefix: "/api", Tags: []string{"BookService"}}
	apiGroup.Use(middlewares.CustomMiddleware)
	apiGroup.Deprecated()
	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "/books",
			Handler:     bookService.List,
			Group:       apiGroup,
			Middlewares: []okapi.Middleware{},
			Options: []okapi.RouteOption{
				okapi.DocSummary("Get Books"),
				okapi.DocDescription("Retrieve a list of books"),
				okapi.DocResponse([]models.Book{}),
				okapi.DocResponse(&models.BookResponse{}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/books/:id",
			Handler: bookService.Get,
			Group:   apiGroup,
			Options: []okapi.RouteOption{
				okapi.DocSummary("Get Book by ID"),
				okapi.DocDescription("Retrieve a book by its ID"),
				okapi.DocPathParam("id", "int", "The ID of the book"),
				okapi.DocResponse(models.Book{}),
				okapi.DocResponse(http.StatusBadRequest, &models.Response[models.Book]{}),
			},
		},
	}
}
func (r *Router) v1BookRoutes() []okapi.RouteDefinition {
	apiGroup := r.group.Group("/books").WithTags([]string{"V1BookService"})

	return []okapi.RouteDefinition{
		{
			Method:      http.MethodGet,
			Path:        "",
			Handler:     bookService.List,
			Middlewares: []okapi.Middleware{middlewares.CustomMiddleware},
			Group:       apiGroup,
			Summary:     "Get Books",
			Description: "Retrieve a list of books",
			Response:    &models.BooksResponse{},
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			Handler:     bookService.Get,
			Middlewares: []okapi.Middleware{middlewares.CustomMiddleware},
			Group:       apiGroup,
			Summary:     "Get Book by ID",
			Description: "Retrieve a book by its ID",
			Options: []okapi.RouteOption{
				okapi.DocPathParam("id", "int", "The ID of the book"),
				okapi.DocResponse(&models.Response[models.Book]{}),
				okapi.DocResponse(http.StatusBadRequest, &models.Response[any]{})},
		},
	}
}

// *************** Auth Routes ****************

func (r *Router) authRoute() okapi.RouteDefinition {
	apiGroup := r.group.Group("/auth").WithTags([]string{"AuthService"})
	apiGroup.Use(middlewares.CustomMiddleware)
	return okapi.RouteDefinition{

		Method:      http.MethodPost,
		Path:        "/login",
		Handler:     authService.Login,
		Group:       apiGroup,
		Summary:     "Login",
		Description: "User login to get a JWT token",
		Request:     &models.AuthRequest{},
		Response:    &models.AuthResponse{},
	}
}

// ************** Authenticated Routes **************

func (r *Router) coreRoutes() []okapi.RouteDefinition {
	coreGroup := r.group.Group("/core").WithTags([]string{"CoreService"})
	// Apply JWT authentication middleware to the admin group
	coreGroup.Use(middlewares.JWTAuth.Middleware)
	coreGroup.Use(middlewares.CustomMiddleware)
	// Enable Bearer token for OpenAPI documentation
	coreGroup.WithSecurity(bearerAuthSecurity)
	return []okapi.RouteDefinition{
		{
			Method:      http.MethodPost,
			Path:        "/whoami",
			Handler:     authService.WhoAmI,
			Group:       coreGroup,
			Summary:     "Whoami",
			Description: "Get the current user's information",
			Response:    &models.Response[models.UserInfo]{},
		},
	}
}

// ***************** Admin Routes *****************

func (r *Router) AdminRoutes() []okapi.RouteDefinition {
	apiGroup := r.group.Group("/admin").WithTags([]string{"AdminService"})
	// Apply JWT authentication middleware to the admin group
	apiGroup.Use(middlewares.AdminJWTAuth.Middleware)
	apiGroup.Use(middlewares.CustomMiddleware)
	apiGroup.WithBearerAuth() // Enable Bearer token for OpenAPI documentation

	return []okapi.RouteDefinition{

		{
			Method:  http.MethodPost,
			Path:    "/books",
			Handler: bookService.Create,
			Group:   apiGroup,
			Options: []okapi.RouteOption{
				okapi.DocSummary("Create Book"),
				okapi.DocDescription("Create a new book"),
				okapi.DocRequestBody(models.Book{}),
			},
			Security: bearerAuthSecurity,
		},

		{
			Method:  http.MethodGet,
			Path:    "/books",
			Handler: bookService.List,
			Group:   apiGroup,
			Options: []okapi.RouteOption{
				okapi.DocSummary("Get Books"),
				okapi.DocDescription("Get books"),
				okapi.DocResponse([]models.Book{}),
			},
			Security: bearerAuthSecurity,
		},
	}
}
