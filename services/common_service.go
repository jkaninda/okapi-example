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
	"fmt"
	"log"
	"net/http"
	"time"

	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"github.com/jkaninda/okapi-example/models"
	"github.com/jkaninda/okapi-example/session"
	okapiws "github.com/jkaninda/okapi-ws"
)

type CommonService struct {
	SessionManager *session.SessionManager
}
type DatastarEvent struct {
	Data map[string]interface{} `json:"data"`
}

// WebSocket upgrades the HTTP connection to WebSocket.
// Config is optional; pass nil to use default settings.
func WebSocket(config *okapiws.WSConfig, c *okapi.Context) (*okapiws.WSConnection, error) {
	upgrader := okapiws.NewWSUpgrader(config)
	return upgrader.Upgrade(c.Response(), c.Request(), nil)
}

// WebSocketWithHeaders upgrades with additional response headers.
func WebSocketWithHeaders(config *okapiws.WSConfig, headers http.Header, c *okapi.Context) (*okapiws.WSConnection, error) {
	upgrader := okapiws.NewWSUpgrader(config)
	return upgrader.Upgrade(c.Response(), c.Request(), headers)
}

// ****************** CommonService *****************

func (cs *CommonService) Home(c *okapi.Context) error {
	return c.Render(http.StatusOK, "home", okapi.M{
		"title":    "Go Okapi basic implementation example ",
		"message":  "Hello from Okapi!",
		"appName":  "Go Okapi Bookstore",
		"headline": "Discover your next great read",
		"books":    books,
	})
}
func (cs *CommonService) Session(c *okapi.Context) error {
	// Generate unique session ID
	sessionID := fmt.Sprintf("session:%s-%s-%s", c.Request().RemoteAddr, goutils.Slug(c.Request().UserAgent()), time.Now().Format("20060102150405"))

	// Add session
	cs.SessionManager.AddSession(sessionID)
	defer cs.SessionManager.RemoveSession(sessionID)

	connectedAt := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ctx := c.Request().Context()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			data := okapi.M{
				"time":              t.Format("15:04:05"),
				"connectedSince":    goutils.FormatDuration(time.Since(connectedAt), 1),
				"activeConnections": cs.SessionManager.GetCount(),
				"totalSessions":     cs.SessionManager.GetTotalCount(),
				"timestamp":         t.Unix(),
			}

			if err := c.SSEvent("message", data); err != nil {
				return err
			}
		}
	}
}

func (cs *CommonService) WhoAmI(c *okapi.Context) error {
	wr := &models.WhoAmIRequest{}
	if err := c.Bind(wr); err != nil {
		return c.AbortBadRequest("Bad request", err)
	}
	if wr.Email == "" {
		logger.Warn("no email found")
	}
	return c.OK(models.WhoAmIResponse{
		Host:   c.Request().Host,
		RealIp: c.RealIP(),
		CurrentUser: models.UserInfo{
			Name:  wr.Name,
			Email: wr.Email,
			Role:  wr.Role,
		},
	})
}

func (cs *CommonService) WebSocketHandle(c *okapi.Context) error {
	token := c.Query("token")
	if len(token) == 0 {
		// validate token here
	}
	ws, err := WebSocket(nil, c)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("error closing WebSocket: %v", err)
		}
	}()

	ws.OnMessage(func(msg *okapiws.WSMessage) {
		log.Printf("[%d] %s", msg.Type, msg.Data)
		// Echo the message back
		_ = ws.Send(msg.Data)
	})

	ws.OnError(func(err error) {
		log.Printf("WebSocket error: %v", err)
	})

	ws.Start()

	// Block until the connection is closed
	<-ws.Context().Done()
	return nil
}
