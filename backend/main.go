package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5"
	"steveucho.com/packages/backend/gameSystem"
	"steveucho.com/packages/backend/gen/sqlQueries"
	"steveucho.com/packages/backend/middleware"
	"steveucho.com/packages/backend/routes"
)

func main() {
	_, exists := os.LookupEnv("DBSTRING")
	if !exists {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}
	dbString := os.Getenv("DBSTRING")

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := pgx.Connect(ctx, dbString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := sqlQueries.New(conn)
	gameOrchestrator := gameSystem.NewGameOrchestrator()
	routes := &routes.App{
		DB:               queries,
		Ctx:              ctx,
		GameOrchestrator: gameOrchestrator,
	}

	router := gin.Default()
	kratos := middleware.NewAuthMiddleware()

	// Websocket group

	websockets := router.Group("/ws")
	websockets.Use(kratos.AuthMiddleware())
	{
		websockets.GET("/lobby/:lobbyID", middleware.GetLobbyContext(gameOrchestrator), routes.JoinWsLobby)
	}

	router.POST("/register/:username", routes.RegisterPlayer)
	player := router.Group("/player")
	player.Use(kratos.AuthMiddleware())
	{
		player.GET("/:username", routes.GetPlayer)
		player.GET("/stats/:id", routes.GetPlayerStats)
		player.GET("/activegame/:id", routes.GetPlayerActiveGame)
	}
	lobbies := router.Group("/lobbies")
	lobbies.Use(kratos.AuthMiddleware())
	{
		lobbies.GET("/all", routes.GetLobbies)
		lobbies.GET("/available/:limit/:offset", routes.GetOpenGames)
		lobbies.POST("/create", routes.CreateBlankLobby)
		lobbies.POST("/create/:name/:buyIn/:maxPlayers", routes.CreateLobby)
	}
	lobby := router.Group("/lobby/:lobbyID")
	lobby.Use(kratos.AuthMiddleware(), middleware.GetLobbyContext(gameOrchestrator))
	{
		lobby.GET("", routes.GetLobby)

		// owner routes
		lobby.PUT("", middleware.ValidateLobbyOwner, routes.UpdateLobby)
		lobby.PUT("/promote-player/:playerID", middleware.ValidateLobbyOwner, routes.UpdateLobbyLeader)
		lobby.DELETE("/remove-player/:playerID", middleware.ValidateLobbyOwner, routes.RemovePlayer)
	}
	game := router.Group("/game/:gameID")
	game.Use(kratos.AuthMiddleware(), middleware.GetGameContext(gameOrchestrator))
	{
		game.GET("/state", routes.GetGameState)
		game.GET("/stock-order-book", routes.GetGameStockOrderBook)
		game.GET("/stock-chart-points", routes.GetStockChartPoints)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
