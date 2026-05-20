package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

type CreateLobbyParams struct {
	Name       string `uri:"name" binding:"required"`
	PlayerID   string `uri:"playerId" binding:"required"`
	BuyIn      int32  `uri:"buyIn" binding:"required"`
	MaxPlayers int32  `uri:"maxPlayers" binding:"required"`
}

func (app *App) CreateLobby(c *gin.Context) {
	var params CreateLobbyParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := app.DB.CreateLobby(app.Ctx, sqlQueries.CreateLobbyParams{
		Name:       params.Name,
		PlayerID:   uuid.MustParse(params.PlayerID),
		BuyIn:      params.BuyIn,
		MaxPlayers: params.MaxPlayers,
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
