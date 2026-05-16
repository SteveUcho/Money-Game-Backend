package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) RegisterPlayer(c *gin.Context) {
	username := c.Param("username")
	id, err := app.DB.CreatePlayer(app.Ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Player created",
		"id":      id,
	})
}
