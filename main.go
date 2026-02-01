package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/config"
	"github.com/jkaninda/okapi-example/routes"
	"github.com/jkaninda/okapi-example/utils"
)

var (
	//go:embed views/*
	Views    embed.FS
	AssetsFS = http.FS(utils.Must(fs.Sub(Views, "views/assets")))
)

func main() {
	// tmpl, _ := okapi.NewTemplateFromDirectory("public/views", ".html", ".tmpl")

	app := okapi.New()
	conf := config.New()
	if err := conf.Initialize(app); err != nil {
		logger.Fatal("Failed to initialize config", "error", err)
	}
	// app.WithRenderer(tmpl)
	app.WithRendererFromFS(Views, "views/*.html")
	route := routes.New(app, conf, AssetsFS)
	route.RegisterRoutes()

	// Start the server
	if err := app.Start(); err != nil {
		panic(err)
	}
}
