package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gameSystem"
)

type GetGameStockOrderBookParams struct {
	GameID string `uri:"gameID" binding:"required"`
}

func (app *App) GetGameStockOrderBook(c *gin.Context) {
	game := c.MustGet("game").(*gameSystem.Game)
	var params GetGameStockOrderBookParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": game.GetStockOrderBook()})
}
