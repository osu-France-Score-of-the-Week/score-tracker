package main

import (
	"log"
	"score-tracker/src/database"
	"score-tracker/src/middlewares"
	"score-tracker/src/models"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.DB.AutoMigrate(&models.Score{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(middlewares.DatabaseMiddleware(database.DB))

	r.Run(":8080")
}
