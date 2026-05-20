package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

type GetGameStateParams struct {
	GameID string `uri:"gameID" binding:"required,uuid"`
}

func (app *App) GetGameState(c *gin.Context) {
	var params GetGameStateParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	gameState, err := app.DB.GetGameState(app.Ctx, uuid.MustParse(params.GameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	playerState, err := app.DB.GetPlayerGameState(app.Ctx, sqlQueries.GetPlayerGameStateParams{
		GameID:   uuid.MustParse(params.GameID),
		PlayerID: uuid.MustParse("params.PlayerID"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"gameState":   gameState,
		"playerState": playerState,
	})
}
