package controllers

import (
	"errors"
	"net/http"
	"score-tracker/middlewares"
	"score-tracker/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetScores(ctx *gin.Context) {
	cursor := ctx.Query("cursor")
	sort := ctx.DefaultQuery("sort", "recent")

	if sort != "recent" && sort != "best" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort (use recent or best)"})
		return
	}

	scoreRepo := repositories.NewScoreRepository(middlewares.GetDB(ctx))
	scoresResponse, err := scoreRepo.GetScores(cursor, sort)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidCursor) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scores"})
		return
	}

	ctx.JSON(http.StatusOK, scoresResponse)
}

func GetScoresByPlayer(ctx *gin.Context) {
	playerID, err := strconv.Atoi(ctx.Param("player_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid player_id"})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	sort := ctx.DefaultQuery("sort", "recent")
	from := ctx.Query("from")
	to := ctx.Query("to")

	scoreRepo := repositories.NewScoreRepository(middlewares.GetDB(ctx))
	scoresResponse, err := scoreRepo.GetScoresByPlayer(playerID, page, sort, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scores"})
		return
	}

	ctx.JSON(http.StatusOK, scoresResponse)
}
