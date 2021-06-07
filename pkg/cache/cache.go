package cache

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/pkg/model"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
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
		log.Fatal("Redis db only allow digit")
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

func (cache *Cache) SetLeaderboard(ctx context.Context, board model.Leaderboard) error {
	value, err := json.Marshal(board)
	if err != nil {
		log.Println(err)
		return err
	}
	return cache.Set(ctx, LeaderboardKey, value, 0).Err()
}
