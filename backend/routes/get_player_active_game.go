package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *App) GetPlayerActiveGame(c *gin.Context) {
	var player struct {
		ID string `uri:"id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	playerActiveGame, err := app.DB.GetPlayerActiveGame(app.Ctx, uuid.MustParse(player.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, playerActiveGame)
}
