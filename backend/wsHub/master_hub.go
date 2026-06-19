package wsHub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetHubRequest struct {
	HubID uuid.UUID
	Res   chan *Hub
}

// ChatHub maintains the set of chat hubs
type MasterHub struct {
	// Registered clients.
	hubs map[uuid.UUID]*Hub

	// Register requests from the clients.
	Register chan uuid.UUID

	// Unregister requests from clients.
	Unregister chan uuid.UUID

	// request to get a hub
	GetHub chan GetHubRequest
}

func NewMasterHub() *MasterHub {
	hub := &MasterHub{
		hubs:       make(map[uuid.UUID]*Hub),
		Register:   make(chan uuid.UUID),
		Unregister: make(chan uuid.UUID),
		GetHub:     make(chan GetHubRequest),
	}
	go hub.Run()
	return hub
}

func (h *MasterHub) Run() {
	for {
		select {
		case hubID := <-h.Register:
			if _, ok := h.hubs[hubID]; !ok {
				hub := NewHub(h)
				go hub.RunManyToMany(hubID)
				h.hubs[hubID] = hub
			}
		case hubID := <-h.Unregister:
			delete(h.hubs, hubID)
		case req := <-h.GetHub:
			if hub, ok := h.hubs[req.HubID]; ok {
				req.Res <- hub
			} else {
				req.Res <- nil
			}
		}
	}
}

type wsHubRequest struct {
	HubID string `uri:"hubID" binding:"required,uuid"`
}

func AddMasterHubContext(masterHub *MasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wsHubRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		masterHub.Register <- uuid.MustParse(req.HubID)
	}
}

func GetMasterHubContext(masterHub *MasterHub) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req wsHubRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		reqChan := make(chan *Hub)
		masterHub.GetHub <- GetHubRequest{HubID: uuid.MustParse(req.HubID), Res: reqChan}
		hub := <-reqChan
		c.Set("hub", hub)
		c.Next()
	}
}
