package gameSystem

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"steveucho.com/packages/backend/wsHub"
)

type GameState struct {
	Symbol         string
	ShareFloat     int
	CurrentPrice   int
	CurrentYear    int
	CurrentQuarter int
	Players        []uuid.UUID // player ids
	StockChart     []ChartPoint
	PlayerCharts   map[uuid.UUID]*PlayerChartData
	PlayerHoldings map[uuid.UUID]*PlayerHolding

	StockOrderBook  map[uuid.UUID]*StockSellOrderEntry
	OptionOrderBook map[uuid.UUID]*OptionSellOrderEntry

	BroadcastHub *wsHub.Hub
}

func NewGame(broadcastHub *wsHub.Hub, symbol string, players []uuid.UUID) *GameState {
	return &GameState{
		Symbol:         symbol,
		ShareFloat:     1000,
		CurrentPrice:   50,
		CurrentYear:    time.Now().Year(),
		CurrentQuarter: 1,
		Players:        players,
		StockChart:     []ChartPoint{},
		PlayerCharts:   make(map[uuid.UUID]*PlayerChartData),
		PlayerHoldings: make(map[uuid.UUID]*PlayerHolding),

		StockOrderBook:  make(map[uuid.UUID]*StockSellOrderEntry),
		OptionOrderBook: make(map[uuid.UUID]*OptionSellOrderEntry),

		BroadcastHub: broadcastHub,
	}
}

func (g *GameState) stockIpoAllocation() {
	var remainingShares = g.ShareFloat
	for _, player := range g.Players {
		holdings, exists := g.PlayerHoldings[player]
		if exists {
			quantity := rand.IntN(remainingShares)
			remainingShares -= quantity
			holdings.Stocks = append(holdings.Stocks, StockOrderFill{
				Quantity: quantity,
				Price:    g.CurrentPrice,
			})
		}
	}
}

func (g *GameState) startGame() {
	g.stockIpoAllocation()
}

func (g *GameState) validateStockSellOrder(playerID uuid.UUID, order []StockSellOrderEntry) bool {
	var totalQuantity int
	for _, o := range order {
		totalQuantity += o.Quantity
	}
	return totalQuantity <= g.PlayerHoldings[playerID].TotalSharesHeld
}

func (g *GameState) validateOptionSellOrder(playerID uuid.UUID, order *OptionSellOrderEntry) bool {
	// TODO: Implement option sell order validation
	return true
}

func (g *GameState) validateStockBuyOrder(playerID uuid.UUID, orderID uuid.UUID) bool {
	// TODO: Implement stock buy order validation
	return true
}

func (g *GameState) validateOptionBuyOrder(playerID uuid.UUID, orderID uuid.UUID) bool {
	// TODO: Implement option buy order validation
	return true
}

func (g *GameState) SubmitTurn(submission TurnSubmission) bool {
	if !g.validateStockSellOrder(submission.PlayerID, submission.StockSellOrders) {
		return false
	}

	for _, order := range submission.OptionSellOrders {
		g.OptionOrderBook[order.ID] = &order
	}
	for _, order := range submission.StockSellOrders {
		g.StockOrderBook[order.ID] = &order
	}
	for _, order := range submission.OptionSellOrders {
		g.OptionOrderBook[order.ID] = &order
	}

	return true
}
