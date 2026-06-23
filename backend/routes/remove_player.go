package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
)

type RemovePlayerRequest struct {
	PlayerID string `uri:"playerID" binding:"required,uuid"`
}

func (app *App) RemovePlayer(c *gin.Context) {
	lobby := c.MustGet("lobby").(*gameSystem.Lobby)
	var req RemovePlayerRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	lobby.RemovePlayer(uuid.MustParse(req.PlayerID))
	c.JSON(200, gin.H{"message": "Player removed"})
}
