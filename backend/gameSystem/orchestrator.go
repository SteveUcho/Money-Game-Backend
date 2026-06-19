package gameSystem

import (
	"github.com/google/uuid"
)

type RegisterLobby struct {
	LobbyID uuid.UUID
	Owner   string
}

type RegisterGame struct {
	LobbyID uuid.UUID
}

type GameOrchestrator struct {
	games   map[uuid.UUID]*GameState
	lobbies map[uuid.UUID]*Lobby

	registerLobby   chan RegisterLobby
	registerGame    chan RegisterGame
	unregisterLobby chan uuid.UUID
	unregisterGame  chan uuid.UUID
}

func (o *GameOrchestrator) Run() {
	for {
		select {
		case lobby := <-o.registerLobby:
			o.lobbies[lobby.LobbyID] = o.createLobby(lobby.Owner)
		case game := <-o.registerGame:
			o.games[game.LobbyID] = o.createGame(game.LobbyID)
		case lobbyID := <-o.unregisterLobby:
			delete(o.lobbies, lobbyID)
		case gameID := <-o.unregisterGame:
			delete(o.games, gameID)
		}
	}
}

func NewGameOrchestrator() *GameOrchestrator {
	gameSystem := &GameOrchestrator{
		games:   make(map[uuid.UUID]*GameState),
		lobbies: make(map[uuid.UUID]*Lobby),

		registerLobby:   make(chan RegisterLobby),
		registerGame:    make(chan RegisterGame),
		unregisterLobby: make(chan uuid.UUID),
		unregisterGame:  make(chan uuid.UUID),
	}
	go gameSystem.Run()
	return gameSystem
}

func (o *GameOrchestrator) createLobby(owner string) *Lobby {
	lobbyID := uuid.New()
	lobby := NewLobby(lobbyID, owner, o)
	o.lobbies[lobbyID] = lobby
	return lobby
}

func (o *GameOrchestrator) createGame(lobbyID uuid.UUID) *GameState {
	gameID := uuid.New()

	game := o.lobbies[lobbyID].StartGame(gameID)
	return game
}

func (o *GameOrchestrator) getGameState(gameID uuid.UUID) *GameState {
	return o.games[gameID]
}

func (o *GameOrchestrator) GetLobby(lobbyID uuid.UUID) (*Lobby, bool) {
	lobby, exists := o.lobbies[lobbyID]
	return lobby, exists
}
