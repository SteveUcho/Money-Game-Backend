package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

func (app *App) GetPlayer(c *gin.Context) {
	username := c.Param("username")
	player, err := app.DB.GetPlayer(app.Ctx, sqlQueries.GetPlayerParams{
		Username: username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, player)
}
