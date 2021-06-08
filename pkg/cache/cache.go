package cache

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/internal/logging"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var (
	logger = logging.ZapL
)

type Cache struct {
	*redis.Client
}

func Init(pool int) *Cache {
	host := os.Getenv("REDIS_HOST")
	username := os.Getenv("REDIS_USER")
	passwd := os.Getenv("REDIS_PASSWD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logger().Fatal("Redis db only allow digit")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Username: username,
		Password: passwd,
		DB:       db,
		PoolSize: pool,
	})
	return &Cache{rdb}
}

func (cache *Cache) SetLeaderboard(ctx context.Context, scoreObjs []repository.ScoreORM) (model.Leaderboard, error) {
	var scores []model.Score
	var leaderboard model.Leaderboard
	for _, score := range scoreObjs {
		scores = append(scores, model.Score{
			ClientID: score.ClientID,
			Score:    score.Score,
		})
	}
	leaderboard = model.Leaderboard{
		TopPlayers: scores,
	}
	value, err := json.Marshal(leaderboard)
	if err != nil {
		logger(ctx).Error(err)
		return leaderboard, err
	}
	return leaderboard, cache.Set(ctx, LeaderboardKey, value, 0).Err()
}

func (cache *Cache) GetLeaderboard(ctx context.Context) (model.Leaderboard, error) {
	var leaderboard model.Leaderboard
	value, err := cache.Get(ctx, LeaderboardKey).Result()

	switch {
	case err == redis.Nil || value == "":
		logger(ctx).Warn("Leaderboard Not found")
		err = ErrNotFound
	case err != nil:
		logger(ctx).Error("Get Leaderboard err: %s\n", err)
	default:
		err = json.Unmarshal([]byte(value), &leaderboard)
		if err != nil {
			logger(ctx).Error(err)
			err = ErrDataCorruption
		}
	}

	return leaderboard, err
}
