package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
	"steveucho.com/packages/backend/models"
)

type UpdateLobbyRequest struct {
	Title      string `json:"title"`
	Symbol     string `json:"symbol"`
	MaxPlayers int    `json:"maxPlayers"`
	BuyIn      int    `json:"buyIn"`
}

func (app *App) UpdateLobby(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	lobby := c.MustGet("lobby").(*gameSystem.Lobby)
	if uuid.MustParse(user.Session.Identity.Id) != lobby.OwnerID {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	var req UpdateLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update lobby fields
	if req.Title != "" {
		lobby.Title = req.Title
	}
	if req.Symbol != "" {
		lobby.Symbol = req.Symbol
	}
	if req.MaxPlayers > 0 {
		lobby.MaxPlayers = req.MaxPlayers
	}
	if req.BuyIn > 0 {
		lobby.BuyIn = req.BuyIn
	}

	c.JSON(200, gin.H{"message": "Lobby updated successfully"})
}
