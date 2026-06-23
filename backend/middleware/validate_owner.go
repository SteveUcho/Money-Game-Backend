package middleware

import (
	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gameSystem"
	"steveucho.com/packages/backend/models"
)

func ValidateLobbyOwner(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	lobby := c.MustGet("lobby").(*gameSystem.Lobby)
	if lobby.OwnerID != user.ID {
		c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
		return
	}
	c.Next()
}
