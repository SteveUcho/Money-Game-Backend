package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetHubRequest struct {
	LobbyID string
	Res     chan *Hub
}

// ChatHub maintains the set of chat hubs
type MasterHub struct {
	// Registered clients.
	hubs map[string]*Hub

	// Register requests from the clients.
	register chan string

	// Unregister requests from clients.
	unregister chan string

	// request to get a hub
	getHub chan GetHubRequest
}

func NewMasterHub() *MasterHub {
	hub := &MasterHub{
		hubs:       make(map[string]*Hub),
		register:   make(chan string),
		unregister: make(chan string),
		getHub:     make(chan GetHubRequest),
	}
	go hub.Run()
	return hub
}

func (h *MasterHub) Run() {
	for {
		select {
		case lobbyID := <-h.register:
			if _, ok := h.hubs[lobbyID]; !ok {
				h.hubs[lobbyID] = NewHub()
			}
		case lobbyID := <-h.unregister:
			delete(h.hubs, lobbyID)
		case req := <-h.getHub:
			if hub, ok := h.hubs[req.LobbyID]; ok {
				req.Res <- hub
			}
		}
	}
}

type wsMasterHubRequest struct {
	HubID string `uri:"hubID" binding:"required"`
}

func AddMasterHubContext(masterHub *MasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wsMasterHubRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		masterHub.register <- req.HubID
		reqChan := make(chan *Hub)
		masterHub.getHub <- GetHubRequest{LobbyID: req.HubID, Res: reqChan}
		hub := <-reqChan
		c.Set("hub", hub)
		c.Set("cleanUpHub", cleanUpHub(masterHub, req.HubID))
		c.Next()
	}
}

func cleanUpHub(chatHubs *MasterHub, lobbyID string) func(*Client) {
	return func(c *Client) {
		if len(c.hub.clients) == 0 {
			chatHubs.unregister <- lobbyID
		}
	}
}
