package gameSystem

import "github.com/google/uuid"

type ChartPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerChartData struct {
	Points []ChartPoint `json:"points"`
	Color  string       `json:"color"`
}

type StockOrderFill struct {
	Quantity int `json:"quantity"`
	Price    int `json:"price"`
}

type OptionOrderFill struct {
	Quantity    int    `json:"quantity"`
	StrikePrice int    `json:"strikePrice"`
	Price       int    `json:"price"`
	Expiration  string `json:"expiration"`
}

type PlayerHolding struct {
	Username        string            `json:"username"`
	TotalSharesHeld int               `json:"totalSharesHeld"`
	Stocks          []StockOrderFill  `json:"stocks"`
	Options         []OptionOrderFill `json:"options"`
	EventCards      []string          `json:"eventCards"`
	Cash            int               `json:"cash"`
}

type StockSellOrderEntry struct {
	ID       uuid.UUID `json:"id"`
	Player   uuid.UUID `json:"player"`
	Quantity int       `json:"quantity"`
	Price    int       `json:"price"`
}

type OptionSellOrderEntry struct {
	ID          uuid.UUID `json:"id"`
	Player      uuid.UUID `json:"player"`
	Quantity    int       `json:"quantity"`
	Price       int       `json:"price"`
	StrikePrice int       `json:"strikePrice"`
	Expiration  string    `json:"expiration"`
}

type TurnSubmission struct {
	PlayerID         uuid.UUID              `json:"playerId"`
	StockSellOrders  []StockSellOrderEntry  `json:"stockSellOrders"`
	OptionSellOrders []OptionSellOrderEntry `json:"optionSellOrders"`
	StockBuyOrders   []uuid.UUID            `json:"stockBuyOrders"`
	OptionBuyOrders  []uuid.UUID            `json:"optionBuyOrders"`
}
