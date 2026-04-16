package routes

import (
	"score-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func GetScoresByPlayer(r *gin.RouterGroup) {
	r.GET("/player/:player_id", controllers.GetScoresByPlayer)
}

func GetScoresByBeatmap(r *gin.RouterGroup) {
	r.GET("/beatmap/:beatmap_id", controllers.GetScoresByBeatmap)
}

func GetScores(r *gin.RouterGroup) {
	r.GET("/", controllers.GetScores)
}

func SetupScoreRoutes(router *gin.RouterGroup) {
	GetScores(router)
	GetScoresByPlayer(router)
	GetScoresByBeatmap(router)
}
