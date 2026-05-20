package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

type GetPlayerParams struct {
	Username string `uri:"username" binding:"required"`
}

func (app *App) GetPlayer(c *gin.Context) {
	var params GetPlayerParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	player, err := app.DB.GetPlayer(app.Ctx, sqlQueries.GetPlayerParams{
		Username: params.Username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, player)
}
