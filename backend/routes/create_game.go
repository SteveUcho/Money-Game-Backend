package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

func (app *App) CreateLobby(c *gin.Context) {
	var game struct {
		Name       string `uri:"name" binding:"required"`
		PlayerID   string `uri:"playerId" binding:"required"`
		BuyIn      int32  `uri:"buyIn" binding:"required"`
		MaxPlayers int32  `uri:"maxPlayers" binding:"required"`
	}
	if err := c.ShouldBindUri(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := app.DB.CreateLobby(app.Ctx, sqlQueries.CreateLobbyParams{
		Name:       game.Name,
		PlayerID:   uuid.MustParse(game.PlayerID),
		BuyIn:      game.BuyIn,
		MaxPlayers: game.MaxPlayers,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Game created",
	})
}
