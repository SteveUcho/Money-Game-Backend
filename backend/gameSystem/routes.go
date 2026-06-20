package gameSystem

import (
	"net/http"
	"slices"

	"maps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/models"
)

func (o *GameOrchestrator) CreateLobby() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		// TODO: Implement lobby creation logic
		theLobby := o.createLobby(user.Traits.Name.First + "'s Lobby")
		c.JSON(http.StatusOK, gin.H{"message": "Lobby created", "lobbyID": theLobby.ID})
	}
}

func (o *GameOrchestrator) GetLobbies() gin.HandlerFunc {
	return func(c *gin.Context) {
		lobbies := slices.Collect(maps.Values(o.lobbies))
		if len(lobbies) == 0 {
			lobbies = []*Lobby{}
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
}

type gameRequest struct {
	GameID string `uri:"gameID" binding:"required,uuid"`
}

func (o *GameOrchestrator) GetGameContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req gameRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		game := o.getGameState(uuid.MustParse(req.GameID))
		c.Set("game", game)
		c.Next()
	}
}
