package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetLobbyRequest struct {
	LobbyID string `uri:"lobbyID" binding:"required,uuid"`
}

type GetLobbyResponsePlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetLobbyResponse struct {
	LobbyID      string                   `json:"lobbyId"`
	Title        string                   `json:"title"`
	OwnerID      string                   `json:"ownerId"`
	Owner        string                   `json:"owner"`
	MaxPlayers   int                      `json:"maxPlayers"`
	Symbol       string                   `json:"symbol"`
	BuyIn        int                      `json:"buyIn"`
	Players      []GetLobbyResponsePlayer `json:"players"`
	PlayersReady []string                 `json:"playersReady"`
}

func (app *App) GetLobby(c *gin.Context) {
	var req GetLobbyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid lobby ID"})
		return
	}
	lobbyID := uuid.MustParse(req.LobbyID)
	lobby, exists := app.GameOrchestrator.GetLobby(lobbyID)
	if !exists {
		c.JSON(404, gin.H{"error": "Lobby not found"})
		return
	}
	players := []GetLobbyResponsePlayer{}
	for id, username := range lobby.Players {
		players = append(players, GetLobbyResponsePlayer{
			ID:   id.String(),
			Name: username,
		})
	}
	c.JSON(200, GetLobbyResponse{
		LobbyID:      lobbyID.String(),
		Title:        lobby.Title,
		OwnerID:      lobby.OwnerID.String(),
		Owner:        lobby.Owner,
		MaxPlayers:   lobby.MaxPlayers,
		Symbol:       lobby.Symbol,
		BuyIn:        lobby.BuyIn,
		Players:      players,
		PlayersReady: []string{},
	})
}
