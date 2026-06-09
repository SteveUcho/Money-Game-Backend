package gameSystem

import (
	"maps"
	"math/rand/v2"
	"slices"

	"github.com/google/uuid"
	"steveucho.com/packages/backend/wsHub"
)

type Lobby struct {
	Title      string
	Owner      string
	Symbol     string
	Players    map[uuid.UUID]string
	BuyIn      int
	MaxPlayers int
	Chat       *wsHub.Hub
	Game       *GameState
}

var sp500Tickers = [10]string{
	"NVDA", "AAPL", "MSFT", "AMZN", "GOOGL",
	"AVGO", "META", "TSLA", "BRK.B", "LLY",
}

func NewLobby(buyIn int, owner string) *Lobby {
	randInt := rand.IntN(len(sp500Tickers))
	selectedTicker := sp500Tickers[randInt]
	chatHub := wsHub.NewHub()
	return &Lobby{
		Title:      owner,
		Owner:      owner,
		Symbol:     selectedTicker,
		Players:    make(map[uuid.UUID]string),
		BuyIn:      buyIn,
		MaxPlayers: 4,
		Chat:       chatHub,
		Game:       nil,
	}
}

func (l *Lobby) StartGame() {
	if l.Game == nil {
		broadcastHub := wsHub.NewHub()
		l.Game = NewGame(broadcastHub, l.Symbol, slices.Collect(maps.Keys(l.Players)))
		l.Game.startGame()
	}
}

func (l *Lobby) AddPlayer(playerID uuid.UUID, playerName string) {
	l.Players[playerID] = playerName
}
