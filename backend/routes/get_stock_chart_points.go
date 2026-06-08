package routes

import "github.com/gin-gonic/gin"

type ChartPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerChartData struct {
	Username string       `json:"username"`
	Points   []ChartPoint `json:"points"`
	Color    string       `json:"color"`
}

type Response struct {
	Players []PlayerChartData `json:"players"`
	Stock   []ChartPoint      `json:"stock"`
}

type GetStockChartPointsParams struct {
	GameID string `uri:"gameID" binding:"required"`
}

func (app *App) GetStockChartPoints(c *gin.Context) {
	var params GetStockChartPointsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	// TODO: Implement logic to fetch stock chart points
	c.JSON(200, Response{
		Players: []PlayerChartData{
			{
				Username: "Player 1",
				Points: []ChartPoint{
					{X: 0, Y: 0},
					{X: 120, Y: 10},
					{X: 280, Y: 20},
					{X: 610, Y: 50},
					{X: 950, Y: 50},
				},
				Color: "red",
			},
			{
				Username: "Player 2",
				Points: []ChartPoint{
					{X: 0, Y: 0},
					{X: 120, Y: 5},
					{X: 280, Y: -300},
					{X: 610, Y: 30},
					{X: 950, Y: -200},
				},
				Color: "blue",
			},
			{
				Username: "Player 3",
				Points: []ChartPoint{
					{X: 0, Y: 0},
					{X: 120, Y: 15},
					{X: 280, Y: -100},
					{X: 610, Y: 40},
					{X: 950, Y: 400},
				},
				Color: "green",
			},
			{
				Username: "Player 4",
				Points: []ChartPoint{
					{X: 0, Y: 0},
					{X: 120, Y: 15},
					{X: 280, Y: 300},
					{X: 610, Y: 40},
					{X: 950, Y: -100},
				},
				Color: "yellow",
			},
		},
		Stock: []ChartPoint{
			{X: 0, Y: 0},
			{X: 120, Y: 0},
			{X: 280, Y: -50},
			{X: 610, Y: 50},
			{X: 950, Y: 50},
		},
	})
}
