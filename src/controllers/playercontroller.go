package controllers

import (
	"net/http"
	"score-tracker/middlewares"
	"score-tracker/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPlayers(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	playerRepo := repositories.NewPlayerRepository(middlewares.GetDB(ctx))
	playersResponse, err := playerRepo.GetPlayers(page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve players"})
		return
	}

	ctx.JSON(http.StatusOK, playersResponse)
}
