package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

func (app *App) GetPlayer(c *gin.Context) {
	user := c.MustGet("user").(*ory.Session)
	player, err := app.DB.GetPlayer(app.Ctx, sqlQueries.GetPlayerParams{
		OryID: user.Identity.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, player)
}
