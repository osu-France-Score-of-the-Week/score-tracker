package repositories

import (
	"fmt"
	"math"
	"score-tracker/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ScoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Create(score *models.Score) error {
	if err := r.db.Create(score).Error; err != nil {
		return err
	}
	return nil
}

func (r *ScoreRepository) GetRecentScores(cursor string) (models.ScoreCursorResponse, error) {
	var scores []models.Score
	query := r.db.Preload("Beatmap").Preload("Beatmap.Beatmapset").Preload("Player").Order("ended_at DESC").Limit(50)
	if cursor != "" {
		query = query.Where("ended_at > ?", cursor)
	}

	if err := query.Find(&scores).Error; err != nil {
		return models.ScoreCursorResponse{}, err
	}

	var nextCursor string
	if len(scores) > 0 {
		nextCursor = strconv.FormatInt(scores[0].EndedAt, 10)
	} else {
		nextCursor = cursor
	}

	return models.ScoreCursorResponse{
		Scores: scores,
		Cursor: nextCursor,
	}, nil
}

func (r *ScoreRepository) GetScoresByPlayer(playerID int, page int, sort string, from string, to string) (models.ScorePageResponse, error) {
	var scores []models.Score

	if page < 1 {
		page = 1
	}

	const limit = 50
	offset := (page - 1) * limit

	query := r.db.
		Preload("Beatmap").
		Preload("Beatmap.Beatmapset").
		Preload("Player").
		Where("player_id = ?", playerID)

	if from != "" {
		fromTime, err := time.Parse("2006-01-02", from)
		if err != nil {
			return models.ScorePageResponse{}, fmt.Errorf("invalid from date: %w", err)
		}
		query = query.Where("ended_at >= ?", fromTime.Unix())
	}
	if to != "" {
		toTime, err := time.Parse("2006-01-02", to)
		if err != nil {
			return models.ScorePageResponse{}, fmt.Errorf("invalid to date: %w", err)
		}
		endOfDay := toTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("ended_at <= ?", endOfDay.Unix())
	}

	var total int64
	if err := query.Model(&models.Score{}).Count(&total).Error; err != nil {
		return models.ScorePageResponse{}, err
	}

	switch sort {
	case "top":
		query = query.Order("pp DESC").Order("ended_at DESC").Order("id DESC")
	case "recent", "":
		fallthrough
	default:
		query = query.Order("ended_at DESC").Order("id DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
		return models.ScorePageResponse{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return models.ScorePageResponse{
		Scores:     scores,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}
