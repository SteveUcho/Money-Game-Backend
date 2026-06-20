package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"steveucho.com/packages/backend/gameSystem"
	"steveucho.com/packages/backend/models"
	"steveucho.com/packages/backend/wsHub"
)

func (app *App) JoinWsLobby(c *gin.Context) {
	lobby := c.MustGet("lobby").(*gameSystem.Lobby)
	user := c.MustGet("user").(*models.User)

	if len(lobby.Players) >= lobby.MaxPlayers {
		c.JSON(400, gin.H{"error": "lobby is full"})
		return
	}

	conn, err := wsHub.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	uID, err := uuid.Parse(user.Session.Identity.Id)
	if err != nil {
		log.Println(err)
		return
	}
	client := gameSystem.NewClient(uID, user.Traits.Name.First, lobby, conn)
	lobby.Join <- client
}
