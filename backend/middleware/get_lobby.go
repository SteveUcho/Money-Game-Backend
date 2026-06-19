package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
)

type getWsLobbyRequest struct {
	LobbyID string `uri:"lobbyID" binding:"required,uuid"`
}

func GetLobbyContext(o *gameSystem.GameOrchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getWsLobbyRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		lobby, exists := o.GetLobby(uuid.MustParse(req.LobbyID))
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
			return
		}
		c.Set("lobby", lobby)
		c.Next()
	}
}
