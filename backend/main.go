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
	"steveucho.com/packages/backend/gen/sqlQueries"
	"steveucho.com/packages/backend/routes"
	"steveucho.com/packages/backend/wsHub"
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
	app := &routes.App{
		DB:  queries,
		Ctx: ctx,
	}

	router := gin.Default()
	kratos := NewMiddleware()

	// Websocket group
	chatHubs := wsHub.NewMasterHub()
	gameHubs := wsHub.NewMasterHub()
	websockets := router.Group("/ws")
	{
		websockets.GET("/lobby/chat/:hubID", kratos.AuthMiddleware(), wsHub.AddMasterHubContext(chatHubs), wsHub.JoinWsHubLobby)
		websockets.GET("/game/events/:hubID", kratos.AuthMiddleware(), wsHub.AddMasterHubContext(gameHubs), wsHub.JoinWsHubLobby)
	}

	router.POST("/register/:username", app.RegisterPlayer)
	player := router.Group("/player")
	{
		player.GET("/:username", kratos.AuthMiddleware(), app.GetPlayer)
		player.GET("/stats/:id", app.GetPlayerStats)
		player.GET("/activegame/:id", app.GetPlayerActiveGame)
	}
	lobby := router.Group("/lobby")
	{
		lobby.POST("/create/:name/:buyIn/:maxPlayers", app.CreateLobby)
		lobby.GET("/available/:limit/:offset", app.GetOpenGames)
	}
	game := router.Group("/game")
	{
		game.GET("/state/:gameID", app.GetGameState)
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
