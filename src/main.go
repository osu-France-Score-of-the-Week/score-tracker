package main

import (
	"log"
	"score-tracker/src/database"
	middlewares "score-tracker/src/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the database schema
	if err := database.DB.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	// Use middlewares
	r.Use(middlewares.DatabaseMiddleware(database.DB))

	r.Run(":8080")
}
