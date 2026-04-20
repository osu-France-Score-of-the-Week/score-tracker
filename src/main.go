package main

import (
	"log"
	"score-tracker/database"
	"score-tracker/jobs"
	"score-tracker/middlewares"
	"score-tracker/models"
	"score-tracker/osuservices"
	"score-tracker/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.DB.AutoMigrate(
		&models.Player{},
		&models.Beatmapset{},
		&models.BeatmapAttributes{},
		&models.Beatmap{},
		&models.Score{},
		&models.Token{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.Use(middlewares.DatabaseMiddleware(database.DB))

	osuSvc := osuservices.NewOsuService(database.DB)

	jobs.CreateJobs(database.DB, osuSvc)

	v1 := r.Group("/api/v1")
	routes.SetupScoreRoutes(v1.Group("/scores"))
	routes.SetupPlayerRoutes(v1.Group("/players"))

	r.Run(":8080")
}
