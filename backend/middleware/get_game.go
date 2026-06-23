package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
)

type gameRequest struct {
	GameID string `uri:"gameID" binding:"required,uuid"`
}

func GetGameContext(o *gameSystem.GameOrchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gameRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		game, exists := o.GetGame(uuid.MustParse(req.GameID))
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.Set("game", game)
		c.Next()
	}
}
