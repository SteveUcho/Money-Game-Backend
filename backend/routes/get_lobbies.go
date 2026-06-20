package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"steveucho.com/packages/backend/gameSystem"
)

func (app *App) GetLobbies(c *gin.Context) {
	lobbies := app.GameOrchestrator.GetLobbiesSlice()
	if len(lobbies) == 0 {
		lobbies = []*gameSystem.Lobby{}
	}
	newSlice := []gin.H{}
	for _, lobby := range lobbies {
		temp := gin.H{
			"id":         lobby.ID,
			"title":      lobby.Title,
			"maxPlayers": lobby.MaxPlayers,
			"players":    len(lobby.Players),
		}
		newSlice = append(newSlice, temp)
	}
	c.JSON(http.StatusOK, gin.H{"lobbies": newSlice})
}
