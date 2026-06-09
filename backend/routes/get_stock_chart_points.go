package routes

import (
	"github.com/gin-gonic/gin"
	game "steveucho.com/packages/backend/gameSystem"
)

type ResponsePlayerChart struct {
	game.PlayerChartData
	Username string `json:"username"`
}

type Response struct {
	Players []ResponsePlayerChart `json:"players"`
	Stock   []game.ChartPoint     `json:"stock"`
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
		Players: []ResponsePlayerChart{
			{
				Username: "Player 1",
				PlayerChartData: game.PlayerChartData{
					Points: []game.ChartPoint{
						{X: 0, Y: 0},
						{X: 120, Y: 10},
						{X: 280, Y: 20},
						{X: 610, Y: 50},
					},
					Color: "red",
				},
			},
			{
				Username: "Player 2",
				PlayerChartData: game.PlayerChartData{
					Points: []game.ChartPoint{
						{X: 0, Y: 0},
						{X: 120, Y: 5},
						{X: 280, Y: -300},
						{X: 610, Y: 30},
						{X: 950, Y: -200},
					},
					Color: "blue",
				},
			},
			{
				Username: "Player 3",
				PlayerChartData: game.PlayerChartData{
					Points: []game.ChartPoint{
						{X: 0, Y: 0},
						{X: 120, Y: 15},
						{X: 280, Y: -100},
						{X: 610, Y: 40},
						{X: 950, Y: 400},
					},
					Color: "purple",
				},
			},
			{
				Username: "Player 4",
				PlayerChartData: game.PlayerChartData{
					Points: []game.ChartPoint{
						{X: 0, Y: 0},
						{X: 120, Y: 15},
						{X: 280, Y: 300},
						{X: 610, Y: 40},
						{X: 950, Y: -100},
					},
					Color: "yellow",
				},
			},
		},
		Stock: []game.ChartPoint{
			{X: 0, Y: 0},
			{X: 120, Y: 0},
			{X: 280, Y: -50},
			{X: 610, Y: 50},
			{X: 950, Y: 50},
		},
	})
}
