package routes

import (
	"score-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func GetRecentScores(r *gin.RouterGroup) {
	r.GET("/recent", controllers.GetScores)
}

func GetScoresByPlayer(r *gin.RouterGroup) {
	r.GET("/player/:player_id", controllers.GetScoresByPlayer)
}

func SetupScoreRoutes(router *gin.RouterGroup) {
	GetRecentScores(router)
	GetScoresByPlayer(router)
}
