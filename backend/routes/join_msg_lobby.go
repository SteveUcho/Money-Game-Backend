package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JoinMsgLobbyRequest struct {
	LobbyID string `uri:"lobbyID" binding:"required"`
}

func AddChatHubContext(chatHubs *ChatHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req JoinMsgLobbyRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chatHubs.register <- req.LobbyID
		reqChan := make(chan *Hub)
		chatHubs.getHub <- GetHubRequest{LobbyID: req.LobbyID, Res: reqChan}
		hub := <-reqChan
		c.Set("hub", hub)
		c.Set("cleanUpHub", cleanUpHub(chatHubs, req.LobbyID))
		c.Next()
	}
}

func cleanUpHub(chatHubs *ChatHub, lobbyID string) func(*Client) {
	return func(c *Client) {
		if len(c.hub.clients) == 0 {
			chatHubs.unregister <- lobbyID
		}
	}
}

func (app *App) JoinMsgLobby(c *gin.Context) {
	hub := c.MustGet("hub").(*Hub)
	cleanUpFunc := c.MustGet("cleanUpHub").(func(*Client))
	var req JoinMsgLobbyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
