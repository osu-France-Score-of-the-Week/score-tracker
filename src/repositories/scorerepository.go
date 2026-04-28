package repositories

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"score-tracker/models"
	"score-tracker/models/responses"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ScoreRepository struct {
	db *gorm.DB
}

var ErrInvalidCursor = errors.New("invalid cursor")

type bestScoresCursor struct {
	Pp      float64 `json:"pp"`
	EndedAt int64   `json:"ended_at"`
	ID      uint    `json:"id"`
}

type recentScoresCursor struct {
	EndedAt int64 `json:"ended_at"`
	ID      uint  `json:"id"`
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

func (r *ScoreRepository) GetScores(cursor string, sort string, start_date string, end_date string) (responses.ScoreCursor, error) {
	const limit = 50

	var scores []responses.ScoreWithAttributes
	query := r.db.Model(&models.Score{}).Preload("Beatmap").Preload("Beatmap.Beatmapset").Preload("Player").Limit(limit)

	switch sort {
	case "best":
		query = query.Limit(limit + 1)

		if start_date != "" {
			startTime, err := time.Parse("2006-01-02", start_date)
			if err != nil {
				return responses.ScoreCursor{}, fmt.Errorf("invalid from date: %w", err)
			}
			query = query.Where("ended_at >= ?", startTime.Unix())
		}

		if end_date != "" {
			endTime, err := time.Parse("2006-01-02", end_date)
			if err != nil {
				return responses.ScoreCursor{}, fmt.Errorf("invalid to date: %w", err)
			}
			endOfDay := endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query = query.Where("ended_at <= ?", endOfDay.Unix())
		}

		query = query.Order("pp DESC").Order("ended_at DESC").Order("id DESC")
		if cursor != "" {
			decodedCursor, err := decodeBestScoresCursor(cursor)
			if err != nil {
				return responses.ScoreCursor{}, err
			}

			query = query.Where(
				"pp < ? OR (pp = ? AND ended_at < ?) OR (pp = ? AND ended_at = ? AND id < ?)",
				decodedCursor.Pp,
				decodedCursor.Pp,
				decodedCursor.EndedAt,
				decodedCursor.Pp,
				decodedCursor.EndedAt,
				decodedCursor.ID,
			)
		}
	default:
		query = query.Order("ended_at DESC").Order("id DESC")
		if cursor != "" {
			decodedCursor, err := decodeRecentScoresCursor(cursor)
			if err != nil {
				return responses.ScoreCursor{}, err
			}

			query = query.Where(
				"ended_at > ? OR (ended_at = ? AND id > ?)",
				decodedCursor.EndedAt,
				decodedCursor.EndedAt,
				decodedCursor.ID,
			)
		}
	}

	if err := query.Find(&scores).Error; err != nil {
		return responses.ScoreCursor{}, err
	}

	hasNextBest := sort == "best" && len(scores) > limit
	if hasNextBest {
		scores = scores[:limit]
	}

	var nextCursor string
	if sort == "recent" && len(scores) > 0 {
		encodedCursor, err := encodeRecentScoresCursor(scores[len(scores)-1])
		if err != nil {
			return responses.ScoreCursor{}, err
		}
		nextCursor = encodedCursor
	} else if sort == "best" && hasNextBest {
		encodedCursor, err := encodeBestScoresCursor(scores[len(scores)-1])
		if err != nil {
			return responses.ScoreCursor{}, err
		}
		nextCursor = encodedCursor
	} else {
		nextCursor = cursor
	}

	// Add beatmap attributes to each score
	for i := range scores {
		var attributes models.BeatmapAttributes
		if err := r.db.Where("beatmap_id = ?", scores[i].BeatmapID).Where("mods = ?", scores[i].Mods).First(&attributes).Error; err == nil {
			scores[i].Attributes = attributes
		}
	}

	return responses.ScoreCursor{
		Scores: scores,
		Cursor: nextCursor,
	}, nil
}

func encodeBestScoresCursor(score responses.ScoreWithAttributes) (string, error) {
	payload := bestScoresCursor{
		Pp:      score.Pp,
		EndedAt: score.EndedAt,
		ID:      score.ID,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func decodeBestScoresCursor(cursor string) (bestScoresCursor, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return bestScoresCursor{}, fmt.Errorf("%w: malformed cursor", ErrInvalidCursor)
	}

	var payload bestScoresCursor
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return bestScoresCursor{}, fmt.Errorf("%w: malformed cursor payload", ErrInvalidCursor)
	}

	if payload.EndedAt <= 0 || payload.ID == 0 {
		return bestScoresCursor{}, fmt.Errorf("%w: missing cursor fields", ErrInvalidCursor)
	}

	return payload, nil
}

func encodeRecentScoresCursor(score responses.ScoreWithAttributes) (string, error) {
	payload := recentScoresCursor{
		EndedAt: score.EndedAt,
		ID:      score.ID,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func decodeRecentScoresCursor(cursor string) (recentScoresCursor, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		// Backward compatibility: allow old numeric cursor (ended_at timestamp).
		endedAt, parseErr := strconv.ParseInt(cursor, 10, 64)
		if parseErr != nil {
			return recentScoresCursor{}, fmt.Errorf("%w: malformed cursor", ErrInvalidCursor)
		}
		if endedAt <= 0 {
			return recentScoresCursor{}, fmt.Errorf("%w: missing cursor fields", ErrInvalidCursor)
		}
		return recentScoresCursor{EndedAt: endedAt, ID: ^uint(0)}, nil
	}

	var payload recentScoresCursor
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return recentScoresCursor{}, fmt.Errorf("%w: malformed cursor payload", ErrInvalidCursor)
	}

	if payload.EndedAt <= 0 || payload.ID == 0 {
		return recentScoresCursor{}, fmt.Errorf("%w: missing cursor fields", ErrInvalidCursor)
	}

	return payload, nil
}

func (r *ScoreRepository) GetScoresByPlayer(playerID int, page int, sort string, from string, to string) (responses.ScorePage, error) {
	var scores []responses.ScoreWithAttributes

	if page < 1 {
		page = 1
	}

	const limit = 50
	offset := (page - 1) * limit

	query := r.db.Model(&models.Score{}).
		Preload("Beatmap").
		Preload("Beatmap.Beatmapset").
		Preload("Player").
		Where("player_id = ?", playerID)

	if from != "" {
		fromTime, err := time.Parse("2006-01-02", from)
		if err != nil {
			return responses.ScorePage{}, fmt.Errorf("invalid from date: %w", err)
		}
		query = query.Where("ended_at >= ?", fromTime.Unix())
	}
	if to != "" {
		toTime, err := time.Parse("2006-01-02", to)
		if err != nil {
			return responses.ScorePage{}, fmt.Errorf("invalid to date: %w", err)
		}
		endOfDay := toTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("ended_at <= ?", endOfDay.Unix())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return responses.ScorePage{}, err
	}

	switch sort {
	case "best":
		query = query.Order("pp DESC").Order("ended_at DESC")
	case "recent", "":
		fallthrough
	default:
		query = query.Order("ended_at DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
		return responses.ScorePage{}, err
	}

	// Load beatmap attributes for each score
	for i := range scores {
		var attributes models.BeatmapAttributes
		if err := r.db.Where("beatmap_id = ?", scores[i].BeatmapID).Where("mods = ?", scores[i].Mods).First(&attributes).Error; err == nil {
			scores[i].Attributes = attributes
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return responses.ScorePage{
		Scores:     scores,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

func (r *ScoreRepository) GetScoresByBeatmap(beatmapID int, page int) (responses.ScorePage, error) {
	var scores []responses.ScoreWithAttributes

	if page < 1 {
		page = 1
	}

	const limit = 50
	offset := (page - 1) * limit

	query := r.db.Model(&models.Score{}).Preload("Beatmap").Preload("Beatmap.Beatmapset").Preload("Player").Where("beatmap_id = ?", beatmapID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return responses.ScorePage{}, err
	}

	if err := query.Order("pp DESC").Order("ended_at DESC").Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
		return responses.ScorePage{}, err
	}

	if len(scores) == 0 {
		return responses.ScorePage{}, gorm.ErrRecordNotFound
	}

	// Load beatmap attributes for each score
	for i := range scores {
		var attributes models.BeatmapAttributes
		if err := r.db.Where("beatmap_id = ?", scores[i].BeatmapID).Where("mods = ?", scores[i].Mods).First(&attributes).Error; err == nil {
			scores[i].Attributes = attributes
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return responses.ScorePage{
		Scores:     scores,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}
