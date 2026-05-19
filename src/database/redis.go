package database

import (
	"context"
	"encoding/json"
	"fmt"
	"errors"
	"score-tracker/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Ctx    context.Context
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	ctx := context.Background()
	return &RedisClient{
		Ctx: ctx,
	}
}

func (r *RedisClient) CachePlayer(player models.Player, ttl time.Duration) error {
	if r.Client == nil {
		return errors.New("redis client is not initialized")
	}

	payload, err := json.Marshal(player)
	if err != nil {
		return err
	}

	return r.Client.Set(r.Ctx, fmt.Sprintf("player:%d", player.ID), payload, ttl).Err()
}

func (r *RedisClient) GetPlayer(playerID uint) (models.Player, bool, error) {
	if r.Client == nil {
		return models.Player{}, false, nil
	}

	raw, err := r.Client.Get(r.Ctx, fmt.Sprintf("player:%d", playerID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return models.Player{}, false, nil
		}
		return models.Player{}, false, err
	}

	var player models.Player
	if err := json.Unmarshal(raw, &player); err != nil {
		return models.Player{}, false, err
	}

	return player, true, nil
}

func (r *RedisClient) ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(r.Ctx).Result()
	if err != nil {
		return nil, err
	}

	r.Client = client
	return client, nil
}
