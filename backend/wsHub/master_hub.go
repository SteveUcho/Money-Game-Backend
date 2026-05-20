package wsHub

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetHubRequest struct {
	HubID string
	Res   chan *Hub
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
		case hubID := <-h.register:
			if _, ok := h.hubs[hubID]; !ok {
				hub := NewHub()
				go hub.RunManyToMany(cleanUpHub(h, hubID))
				h.hubs[hubID] = hub
			}
		case hubID := <-h.unregister:
			delete(h.hubs, hubID)
		case req := <-h.getHub:
			if hub, ok := h.hubs[req.HubID]; ok {
				req.Res <- hub
			}
		}
	}
}

func cleanUpHub(masterHub *MasterHub, hubID string) func() {
	return func() {
		masterHub.unregister <- hubID
	}
}

type wsHubRequest struct {
	HubID string `uri:"hubID" binding:"required"`
}

func AddMasterHubContext(masterHub *MasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wsHubRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		masterHub.register <- req.HubID
		reqChan := make(chan *Hub)
		masterHub.getHub <- GetHubRequest{HubID: req.HubID, Res: reqChan}
		hub := <-reqChan
		c.Set("hub", hub)
		c.Next()
	}
}
