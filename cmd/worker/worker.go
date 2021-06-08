package main

import (
	"context"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/repository"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

var (
	DB  *repository.Repository
	RDB *cache.Cache
)

type ResetLeaderboard struct {
	DB  *repository.Repository
	RDB *cache.Cache
}

func (j ResetLeaderboard) Run() {
	ctx := context.Background()
	topTen, err := j.DB.ListTopScores(ctx)
	_, err = j.RDB.SetLeaderboard(ctx, topTen)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var wg sync.WaitGroup
	DB = repository.Init(20, 1)
	RDB = cache.Init(10)
	log.Println("worker start")
	c := cron.New()
	c.AddJob("@every 10m", ResetLeaderboard{DB, RDB})
	wg.Add(1)
	c.Start()
	wg.Wait()
}
