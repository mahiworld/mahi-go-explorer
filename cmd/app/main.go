package main

import (
	"log"
	"mahi-go-explorer/internal/api/handlers"
	"mahi-go-explorer/internal/config"
	"mahi-go-explorer/internal/store"
	userpkg "mahi-go-explorer/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//create a new gin app
	app := gin.New()

	//enable cors
	// corsConfig := cors.DefaultConfig()
	// app.Use(cors.New(corsConfig))

	// define a ping route
	app.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    "pong",
		})
	})

	//connect to the database
	db, err := store.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	cc := config.CreateCollection()
	log.Println("Connected to database")

	//create services
	userService := userpkg.NewService(db, cc)

	//register routes
	handlers.RegisterRoutes(
		app,
		userService,
	)

	//Ensure admin user exists
	if err = userService.EnsureAdminUserExists(); err != nil {
		log.Fatalf("Error ensuring admin user exists: %v", err)
	}

	app.Run(":8080")
	log.Println("Server started on port 8080")
}
