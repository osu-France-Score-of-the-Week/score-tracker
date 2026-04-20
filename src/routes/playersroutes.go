package routes

import (
	"score-tracker/controllers"

	"github.com/gin-gonic/gin"
)

func GetPlayers(r *gin.RouterGroup) {
	r.GET("/", controllers.GetPlayers)
}

func SetupPlayerRoutes(router *gin.RouterGroup) {
	GetPlayers(router)
}
