package gameSystem

import (
	"encoding/json"
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"

	"github.com/google/uuid"
)

type Lobby struct {
	ID         uuid.UUID
	Title      string
	OwnerID    uuid.UUID
	Owner      string
	Symbol     string
	Players    map[uuid.UUID]string
	BuyIn      int
	MaxPlayers int
	Game       *Game

	Clients      map[*Client]uuid.UUID
	Join         chan *Client
	Leave        chan *Client
	Broadcast    chan ClientBroadcast
	orchestrator *GameOrchestrator
}

var sp500Tickers = [10]string{
	"nvda", "aapl", "msft", "amzn", "googl",
	"avgo", "meta", "tsla", "lly",
}

func NewLobby(lobbyID uuid.UUID, ownerID uuid.UUID, ownerUsername string, orchestrator *GameOrchestrator) *Lobby {
	randInt := rand.IntN(len(sp500Tickers))
	selectedTicker := sp500Tickers[randInt]
	theLobby := &Lobby{
		ID:         lobbyID,
		Title:      ownerUsername,
		OwnerID:    ownerID,
		Owner:      ownerUsername,
		Symbol:     selectedTicker,
		Players:    make(map[uuid.UUID]string),
		BuyIn:      0,
		MaxPlayers: 4,
		Game:       nil,

		Clients:      make(map[*Client]uuid.UUID),
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		Broadcast:    make(chan ClientBroadcast),
		orchestrator: orchestrator,
	}
	go theLobby.Run()
	return theLobby
}

func (l *Lobby) Run() {
	for {
		select {
		case client := <-l.Join:
			l.addPlayer(client)
		case client := <-l.Leave:
			l.removeClient(client)
		case broadcast := <-l.Broadcast:
			for client := range l.Clients {
				select {
				case client.Send <- broadcast:
				default:
					l.removeClient(client)
				}
			}
		}
	}
}

func (l *Lobby) StartGame(gameID uuid.UUID) *Game {
	if l.Game == nil {
		l.Game = NewGame(gameID, l.Symbol, slices.Collect(maps.Keys(l.Players)))
		l.Game.startGame()
	}
	return l.Game
}

type SystemAction struct {
	Action   string    `json:"action"`
	PlayerID uuid.UUID `json:"playerId"`
	Username string    `json:"username"`
}

type ChatMessage struct {
	Message string `json:"message"`
}

func (l *Lobby) addPlayer(client *Client) {
	oldPlayerCount := len(l.Players)
	l.Players[client.ID] = client.username
	newPlayerCount := len(l.Players)
	l.Clients[client] = client.ID

	if oldPlayerCount != newPlayerCount {
		go func() { // without goroutine, it will infinite block
			// player joined
			playerJoinBroadcast := ChatMessage{
				Message: fmt.Sprintf("%s joined the lobby", client.username),
			}
			jsonData, err := json.Marshal(playerJoinBroadcast)
			if err != nil {
				return
			}
			l.Broadcast <- ClientBroadcast{
				Type:      "chat",
				MessageID: uuid.New(),
				PlayerID:  uuid.New(),
				Username:  "System",
				Data:      jsonData,
			}
			newMutationBroadcast := map[string][]GetLobbyResponsePlayer{
				"players": l.GetPlayers(),
			}
			jsonData, err = json.Marshal(newMutationBroadcast)
			if err != nil {
				return
			}
			l.Broadcast <- ClientBroadcast{
				Type:      "lobby.update",
				MessageID: uuid.New(),
				PlayerID:  uuid.New(),
				Username:  "System",
				Data:      jsonData,
			}
		}()
	}
}

type GetLobbyResponsePlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (l *Lobby) GetPlayers() []GetLobbyResponsePlayer {
	players := []GetLobbyResponsePlayer{}
	for id, username := range l.Players {
		players = append(players, GetLobbyResponsePlayer{
			ID:   id.String(),
			Name: username,
		})
	}
	return players
}

// player may have multiple connections, so we need to remove all of them
func (l *Lobby) RemovePlayer(playerID uuid.UUID) {
	clientsToDelete := make([]*Client, 0)
	for client, id := range l.Clients {
		if id == playerID {
			clientsToDelete = append(clientsToDelete, client)
		}
	}
	for _, client := range clientsToDelete {
		l.Leave <- client
	}
}

func (l *Lobby) removeClient(client *Client) {
	if _, exists := l.Clients[client]; !exists {
		return
	}
	delete(l.Clients, client)
	close(client.Send)

	if len(l.Clients) == 0 {
		l.closeLobby()
		return
	}

	newPlayers := make(map[uuid.UUID]string)
	for c := range l.Clients {
		newPlayers[c.ID] = c.username
	}
	if len(newPlayers) != len(l.Players) {
		l.Players = newPlayers
		if client.ID == l.OwnerID {
			// pick a new owner
			for c := range l.Clients {
				l.OwnerID = c.ID
				l.Owner = c.username
				break
			}
		}

		playerLeftBroadcast := ChatMessage{
			Message: fmt.Sprintf("%s left the lobby", client.username),
		}
		jsonData, err := json.Marshal(playerLeftBroadcast)
		if err != nil {
			return
		}
		go func() { // without goroutine, it will infinite block
			l.Broadcast <- ClientBroadcast{
				Type:      "chat",
				MessageID: uuid.New(),
				PlayerID:  uuid.New(),
				Username:  "System",
				Data:      jsonData,
			}

			newMutationBroadcast := map[string][]GetLobbyResponsePlayer{
				"players": l.GetPlayers(),
			}
			jsonData, err := json.Marshal(newMutationBroadcast)
			if err != nil {
				return
			}
			l.Broadcast <- ClientBroadcast{
				Type:      "lobby.update",
				MessageID: uuid.New(),
				PlayerID:  uuid.New(),
				Username:  "System",
				Data:      jsonData,
			}
		}()
	}
}

func (l *Lobby) UpdateLeader(playerID uuid.UUID) {
	l.OwnerID = playerID
	l.Owner = l.Players[playerID]

	newChatBroadcast := ChatMessage{
		Message: fmt.Sprintf("%s is now the lobby leader", l.Owner),
	}
	jsonData, err := json.Marshal(newChatBroadcast)
	if err != nil {
		return
	}
	l.Broadcast <- ClientBroadcast{
		Type:      "chat",
		MessageID: uuid.New(),
		PlayerID:  uuid.New(),
		Username:  "System",
		Data:      jsonData,
	}
	newMutationBroadcast := map[string]string{
		"owner":   l.Owner,
		"ownerID": l.OwnerID.String(),
	}
	jsonData, err = json.Marshal(newMutationBroadcast)
	if err != nil {
		return
	}
	l.Broadcast <- ClientBroadcast{
		Type:      "lobby.update",
		MessageID: uuid.New(),
		PlayerID:  uuid.New(),
		Username:  "System",
		Data:      jsonData,
	}
}

func (l *Lobby) closeLobby() {
	l.orchestrator.unregisterLobby <- l.ID
	for client := range l.Clients {
		close(client.Send)
	}
	if l.Game != nil {
		l.Game.EndGame()
	}
}
