package cache

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/pkg/model"
	"fmt"
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

func (cache *Cache) GetLeaderboard(ctx context.Context) (model.Leaderboard, error) {
	var leaderboard model.Leaderboard
	value, err := cache.Get(ctx, LeaderboardKey).Result()

	switch {
	case err == redis.Nil || value == "":
		fmt.Println("Leaderboard Not found")
		err = ErrNotFound
	case err != nil:
		fmt.Printf("Get Leaderboard err: %s\n", err)
	default:
		err = json.Unmarshal([]byte(value), &leaderboard)
		if err != nil {
			fmt.Println(err)
			err = ErrDataCorruption
		}
	}

	return leaderboard, err
}
