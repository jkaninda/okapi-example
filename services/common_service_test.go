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

package services

import (
	"html/template"
	"io"
	"testing"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/session"
	"github.com/jkaninda/okapi/okapitest"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c *okapi.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var commonService = &CommonService{SessionManager: session.New()}

func TestCommonAPI(t *testing.T) {
	// Setup test server
	server := okapi.NewTestServer(t)
	server.With().WithRenderer(&Template{templates: template.Must(template.ParseGlob("../views/*.html"))})

	server.Get("/", commonService.Home)

	// Create reusable client
	client := okapitest.NewClient(t, server.BaseURL)

	client.GET("/").ExpectStatusOK().ExpectBodyContains("Go Okapi Bookstore")

}
