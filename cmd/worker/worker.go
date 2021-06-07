package main

import (
	"context"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	DB *repository.Repository
	Cache *cache.Cache
)

func ResetLeaderboard() {
	var scores []model.Score
	ctx := context.Background()
	topTen, err := DB.ListTopScores(ctx, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, score := range topTen {
		scores = append(scores, model.Score{
			ClientID: score.ClientID,
			Score: score.Score,
		})
	}
	leaderboard := model.Leaderboard{
		TopPlayers: scores,
	}
	err = Cache.SetLeaderboard(ctx, leaderboard)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var wg sync.WaitGroup
	DB = repository.Init(20, 1)
	Cache = cache.Init(10)
	c := cron.New()
	c.AddFunc("@every 10m", ResetLeaderboard)
	wg.Wait()
}
