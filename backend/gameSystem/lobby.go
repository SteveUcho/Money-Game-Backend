package gameSystem

import (
	"encoding/json"
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
	Game       *GameState

	Clients      map[*Client]bool
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

		Clients:      make(map[*Client]bool),
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
			l.removePlayer(client)
		case broadcast := <-l.Broadcast:
			for client := range l.Clients {
				select {
				case client.Send <- broadcast:
				default:
					close(client.Send)
					delete(l.Clients, client)
				}
			}
		}
	}
}

func (l *Lobby) StartGame(gameID uuid.UUID) *GameState {
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

func (l *Lobby) addPlayer(client *Client) {
	l.Players[client.ID] = client.username
	l.Clients[client] = true

	playerJoinBroadcast := SystemAction{
		Action:   "player_joined",
		PlayerID: client.ID,
		Username: client.username,
	}
	jsonData, err := json.Marshal(playerJoinBroadcast)
	if err != nil {
		return
	}
	go func() { // without goroutine, it will infinite block
		l.Broadcast <- ClientBroadcast{
			Type:     "system",
			PlayerID: uuid.New(),
			Username: "System",
			Data:     jsonData,
		}
	}()
}

func (l *Lobby) removePlayer(client *Client) {
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

		playerLeftBroadcast := SystemAction{
			Action:   "player_left",
			PlayerID: client.ID,
			Username: client.username,
		}
		jsonData, err := json.Marshal(playerLeftBroadcast)
		if err != nil {
			return
		}
		go func() { // without goroutine, it will infinite block
			l.Broadcast <- ClientBroadcast{
				Type:     "system",
				PlayerID: uuid.New(),
				Username: "System",
				Data:     jsonData,
			}
		}()
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
