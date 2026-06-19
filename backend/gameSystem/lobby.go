package gameSystem

import (
	"maps"
	"math/rand/v2"
	"slices"

	"github.com/google/uuid"
)

type Lobby struct {
	ID         uuid.UUID
	Title      string
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
	"NVDA", "AAPL", "MSFT", "AMZN", "GOOGL",
	"AVGO", "META", "TSLA", "BRK.B", "LLY",
}

func NewLobby(lobbyID uuid.UUID, owner string, orchestrator *GameOrchestrator) *Lobby {
	randInt := rand.IntN(len(sp500Tickers))
	selectedTicker := sp500Tickers[randInt]
	theLobby := &Lobby{
		ID:         lobbyID,
		Title:      owner,
		Owner:      owner,
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
			l.Clients[client] = true
		case client := <-l.Leave:
			if _, ok := l.Clients[client]; ok {
				delete(l.Clients, client)
				close(client.Send)
			}
			if len(l.Clients) == 0 {
				l.closeLobby()
			}
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

func (l *Lobby) AddPlayer(playerID uuid.UUID, playerName string) {
	l.Players[playerID] = playerName
}

// TODO: close lobby when all players disconnect
func (l *Lobby) closeLobby() {
	l.orchestrator.unregisterLobby <- l.ID
	if l.Game != nil {
		l.Game.EndGame()
	}
}
