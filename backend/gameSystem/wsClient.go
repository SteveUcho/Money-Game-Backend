// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gameSystem

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ClientBroadcast struct {
	Type      string
	MessageID uuid.UUID
	PlayerID  uuid.UUID
	Username  string
	Data      []byte
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID       uuid.UUID
	username string
	lobby    *Lobby

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan ClientBroadcast
}

func NewClient(id uuid.UUID, username string, lobby *Lobby, conn *websocket.Conn) *Client {
	client := &Client{
		ID:       id,
		username: username,
		lobby:    lobby,
		conn:     conn,
		Send:     make(chan ClientBroadcast, 256),
	}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return client
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.lobby.Leave <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		c.lobby.Broadcast <- ClientBroadcast{
			Type:      "chat",
			MessageID: uuid.New(),
			PlayerID:  c.ID,
			Username:  c.username,
			Data:      message,
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case broadcast, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			var msg map[string]any
			if err := json.Unmarshal(broadcast.Data, &msg); err != nil {
				log.Printf("error unmarshalling message: %v", err)
				return
			}
			test := map[string]any{
				"type":      broadcast.Type,
				"playerId":  broadcast.PlayerID,
				"messageId": broadcast.MessageID.String(),
				"username":  broadcast.Username,
				"data":      msg,
			}
			message, err := json.Marshal(test)
			if err != nil {
				log.Printf("error marshalling to JSON: %v", err)
				return
			}
			w.Write(message)

			// Add queued ws messages to the current websocket message.
			n := len(c.Send)
			for range n {
				w.Write(newline)
				broadcast := <-c.Send
				var msg map[string]any
				if err := json.Unmarshal(broadcast.Data, &msg); err != nil {
					log.Printf("error unmarshalling message: %v", err)
					return
				}
				temp := map[string]any{
					"type":      broadcast.Type,
					"playerId":  broadcast.PlayerID,
					"messageId": broadcast.MessageID.String(),
					"username":  broadcast.Username,
					"data":      msg,
				}
				message, err := json.Marshal(temp)
				if err != nil {
					log.Printf("error marshalling to JSON: %v", err)
					return
				}
				w.Write(message)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
