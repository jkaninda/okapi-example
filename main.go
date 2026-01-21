package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/config"
	"github.com/jkaninda/okapi-example/routes"
)

//go:embed views/*
var Views embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c *okapi.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func NewTemplate() *Template {
	tmpl := template.Must(template.ParseFS(Views, "views/*.html"))
	return &Template{templates: tmpl}
}
func main() {
	app := okapi.New()
	conf := config.New()
	if err := conf.Initialize(app); err != nil {
		logger.Fatal("Failed to initialize config", "error", err)
	}
	app.WithRenderer(NewTemplate())
	route := routes.New(app, conf)
	route.RegisterRoutes()

	// Start the server
	if err := app.Start(); err != nil {
		panic(err)
	}
}
