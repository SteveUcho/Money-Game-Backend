// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wsHub

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/models"
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
	// Master hub that manages all hubs
	master *MasterHub

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan ClientBroadcast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub(master *MasterHub) *Hub {
	return &Hub{
		master:     master,
		broadcast:  make(chan ClientBroadcast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) RunManyToMany(hubID uuid.UUID) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			if len(h.clients) == 0 {
				h.master.Unregister <- hubID
				return
			}
		case broadcast := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- broadcast:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func JoinWsHubLobby(c *gin.Context) {
	hub := c.MustGet("hub").(*Hub)
	user := c.MustGet("user").(*models.User)

	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(user.Session.Identity.Id, user.Traits.Name.First, hub, conn)
	hub.register <- client
}
