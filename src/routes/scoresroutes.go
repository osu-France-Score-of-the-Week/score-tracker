package routes

import (
	"score-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func GetScoresByPlayer(r *gin.RouterGroup) {
	r.GET("/player/:player_id", controllers.GetScoresByPlayer)
}

func GetScores(r *gin.RouterGroup) {
	r.GET("/", controllers.GetScores)
}

func SetupScoreRoutes(router *gin.RouterGroup) {
	GetScores(router)
	GetScoresByPlayer(router)
}
