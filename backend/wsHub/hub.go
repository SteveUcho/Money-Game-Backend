// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wsHub

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ory "github.com/ory/kratos-client-go"
)

type ClientBroadcast struct {
	ClientID  string
	Username  string
	MessageID uuid.UUID
	Message   []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan ClientBroadcast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan ClientBroadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) RunManyToMany(cleanUpFunc func()) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 0 {
				cleanUpFunc()
				return
			}
		case broadcast := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- broadcast:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func JoinWsHubLobby(c *gin.Context) {
	hub := c.MustGet("hub").(*Hub)
	user := c.MustGet("user").(*ory.Session)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(user.Identity.Id, user.Identity.Traits.(map[string]any)["name"].(map[string]any)["first"].(string), hub, conn)
	hub.register <- client
}
