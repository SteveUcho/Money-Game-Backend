package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetPlayerStatsParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (app *App) GetPlayerStats(c *gin.Context) {
	var params GetPlayerStatsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	playerStats, err := app.DB.GetPlayerStats(app.Ctx, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, playerStats)
}
