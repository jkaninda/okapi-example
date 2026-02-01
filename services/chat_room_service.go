package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	okapiws "github.com/jkaninda/okapi-ws"
)

type ChatRoomService struct{}

type ChatMessage struct {
	// Type event type, "join", "leave", "message", "system"
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	UserCount int       `json:"userCount,omitempty"`
}

type Client struct {
	ws       *okapiws.WSConnection
	username string
	room     *ChatRoom
}

type ChatRoom struct {
	clients    map[*Client]bool
	broadcast  chan ChatMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var chatRoom = &ChatRoom{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan ChatMessage, 100),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// init starts the chat room
func init() {
	go chatRoom.run()
}

func (room *ChatRoom) run() {
	for {
		select {
		case client := <-room.register:
			room.mu.Lock()
			room.clients[client] = true
			userCount := len(room.clients)
			room.mu.Unlock()

			// Broadcast join message
			room.broadcast <- ChatMessage{
				Type:      "join",
				Username:  client.username,
				Message:   fmt.Sprintf("%s joined the chat", client.username),
				Timestamp: time.Now(),
				UserCount: userCount,
			}

		case client := <-room.unregister:
			room.mu.Lock()
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				userCount := len(room.clients)
				room.mu.Unlock()

				// Broadcast leave message
				room.broadcast <- ChatMessage{
					Type:      "leave",
					Username:  client.username,
					Message:   fmt.Sprintf("%s left the chat", client.username),
					Timestamp: time.Now(),
					UserCount: userCount,
				}
			} else {
				room.mu.Unlock()
			}

		case message := <-room.broadcast:
			room.mu.RLock()
			for client := range room.clients {
				if err := client.sendMessage(message); err != nil {
					log.Printf("Error sending to client %s: %v", client.username, err)
				}
			}
			room.mu.RUnlock()
		}
	}
}

// getUserCount returns the current number of connected clients
func (room *ChatRoom) getUserCount() int {
	room.mu.RLock()
	defer room.mu.RUnlock()
	return len(room.clients)
}

// sendMessage sends a message to this client
func (c *Client) sendMessage(msg ChatMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return c.ws.Send(data)
}

// WebSocketHandle handles WebSocket connections for the chat room
func (cs *ChatRoomService) WebSocketHandle(c *okapi.Context) error {
	// Get username from query parameter
	username := c.Query("username")
	if len(username) == 0 {
		return c.AbortBadRequest("Bad Request", fmt.Errorf("missing username"))
	}

	// Upgrade to WebSocket
	ws, err := WebSocket(nil, c)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("error closing WebSocket: %v", err)
		}
	}()

	// Create client
	client := &Client{
		ws:       ws,
		username: username,
		room:     chatRoom,
	}

	// Register client
	chatRoom.register <- client

	// Cleanup on disconnect
	defer func() {
		chatRoom.unregister <- client
		if err := ws.Close(); err != nil {
			log.Printf("Error closing WebSocket for %s: %v", username, err)
		}
	}()

	// Send welcome message
	welcomeMsg := ChatMessage{
		Type:      "system",
		Message:   "Welcome to the chat room! Ask a question to start the conversation.",
		Timestamp: time.Now(),
		UserCount: chatRoom.getUserCount(),
	}
	if err := client.sendMessage(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
	}
	ws.OnMessage(func(msg *okapiws.WSMessage) {
		var chatMsg ChatMessage
		if err := json.Unmarshal(msg.Data, &chatMsg); err != nil {
			log.Printf("Error unmarshaling message from %s: %v", username, err)
			return
		}

		if chatMsg.Type != "message" {
			return
		}

		chatMsg.Username = username
		chatMsg.Timestamp = time.Now()
		chatMsg.UserCount = chatRoom.getUserCount()

		// Broadcast to all clients
		chatRoom.broadcast <- chatMsg
	})

	ws.OnError(func(err error) {
		logger.Error("WebSocket error", "username", username)
	})

	ws.OnClose(func() {
		logger.Info("WebSocket closed", "username", username)
	})

	ws.Start()

	<-ws.Context().Done()

	return nil
}

// ChatPage renders the chat page
func (cs *ChatRoomService) ChatPage(c *okapi.Context) error {
	return c.Render(200, "chat", okapi.M{
		"title":    "Okapi simple Chat Room",
		"appName":  "Live Chat",
		"headline": "Connect and chat in real-time",
	})
}
