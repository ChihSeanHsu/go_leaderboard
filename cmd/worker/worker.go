package main

import (
	"context"
	"example.com/leaderboard/internal/logging"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/repository"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	DB     *repository.Repository
	Cache  *cache.Cache
	logger = logging.ZapL
)

type ResetLeaderboard struct {
	DB    *repository.Repository
	Cache *cache.Cache
}

func (j ResetLeaderboard) Run() {
	ctx := context.WithValue(context.Background(), "TraceID", uuid.New().String())
	topTen, err := j.DB.ListTopScores(ctx)
	_, err = j.Cache.SetLeaderboard(ctx, topTen)
	if err != nil {
		logger(ctx).Error(err)
	}
}

func main() {
	var wg sync.WaitGroup
	logging.InitLogging()
	DB = repository.Init(20, 1)
	Cache = cache.Init(10)

	c := cron.New()
	c.AddJob("@every 10m", ResetLeaderboard{DB, Cache})
	wg.Add(1)

	logger().Info("worker start")
	c.Start()
	wg.Wait()
}
