package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
)

type UpdateLobbyLeaderRequest struct {
	PlayerID string `uri:"playerID" binding:"required,uuid"`
}

func (app *App) UpdateLobbyLeader(c *gin.Context) {
	var req UpdateLobbyLeaderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	lobby := c.MustGet("lobby").(*gameSystem.Lobby)
	playerID := uuid.MustParse(req.PlayerID)

	if _, exists := lobby.Players[playerID]; !exists {
		c.JSON(404, gin.H{"error": "Player not found"})
		return
	}

	lobby.UpdateLeader(playerID)
	c.JSON(200, gin.H{"message": "Lobby leader updated"})
}
