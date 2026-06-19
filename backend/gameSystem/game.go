package gameSystem

import (
	"maps"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/google/uuid"
	"steveucho.com/packages/backend/wsHub"
)

type GameState struct {
	ID             uuid.UUID
	Symbol         string
	ShareFloat     int
	CurrentPrice   int
	CurrentYear    int
	CurrentQuarter int
	Players        []uuid.UUID // player ids
	StockChart     []ChartPoint
	PlayerCharts   map[uuid.UUID]*PlayerChartData
	PlayerHoldings map[uuid.UUID]*PlayerHolding

	StockOrderBook  map[uuid.UUID]*StockSellOrderEntry  // order id -> order
	OptionOrderBook map[uuid.UUID]*OptionSellOrderEntry // order id -> order

	BroadcastHub *wsHub.Hub
	lobby        *Lobby
}

func NewGame(ID uuid.UUID, symbol string, players []uuid.UUID) *GameState {
	return &GameState{
		ID:             ID,
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

func (g *GameState) validateStockBuyOrder(playerID uuid.UUID, orderID uuid.UUID) bool {
	// TODO: Implement stock buy order validation
	return true
}

func (g *GameState) resolveStockSellOrder(orderID uuid.UUID) {
	order, ok := g.StockOrderBook[orderID]
	if !ok {
		return
	}

	newStocks := slices.Clone(g.PlayerHoldings[order.Player].Stocks)
	slices.SortFunc(g.PlayerHoldings[order.Player].Stocks, func(a, b StockOrderFill) int {
		return a.Price - b.Price
	})
	remainingQuantity := order.Quantity
	deleteIndexes := make([]int, len(newStocks))
	index := 0
	for remainingQuantity > 0 {
		stock := newStocks[index]
		if stock.Quantity <= remainingQuantity {
			deleteIndexes = append(deleteIndexes, index)
		} else {
			stock.Quantity -= remainingQuantity
			newStocks[index] = stock
		}
		index++
	}
	index = 1
	for _, delIndex := range deleteIndexes { // move all to be deleted stocks to the end of the slice
		newStocks[delIndex] = newStocks[len(newStocks)-index]
		index++
	}
	newStocks = slices.Delete(newStocks, len(newStocks)-len(deleteIndexes), len(newStocks)) // delete the stocks at the end of the slice

	g.PlayerHoldings[order.Player].Stocks = newStocks
	g.PlayerHoldings[order.Player].TotalSharesHeld -= order.Quantity
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

func (g *GameState) GetStockOrderBook() []*StockSellOrderEntry {
	return slices.Collect(maps.Values(g.StockOrderBook))
}

func (g *GameState) EndGame() {
	g.lobby.orchestrator.unregisterGame <- g.ID
	g.lobby.Game = nil
}
