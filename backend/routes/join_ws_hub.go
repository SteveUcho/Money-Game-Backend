package routes

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (app *App) JoinWsHubLobby(c *gin.Context) {
	hub := c.MustGet("hub").(*Hub)
	cleanUpFunc := c.MustGet("cleanUpHub").(func(*Client))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(cleanUpFunc)
}
