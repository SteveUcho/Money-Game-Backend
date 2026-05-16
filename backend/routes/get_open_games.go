package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gen/sqlQueries"
)

type GetOpenGamesRequest struct {
	Limit  int32 `uri:"limit" binding:"required"`
	Offset int32 `uri:"offset" binding:"required"`
}

func (app *App) GetOpenGames(c *gin.Context) {
	var req GetOpenGamesRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	openLobbies, err := app.DB.GetOpenLobbies(app.Ctx, sqlQueries.GetOpenLobbiesParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, openLobbies)
}
