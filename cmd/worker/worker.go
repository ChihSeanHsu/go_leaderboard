package main

import (
	"context"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
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
	var scores []model.Score
	ctx := context.Background()
	topTen, err := j.DB.ListTopScores(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	for _, score := range topTen {
		scores = append(scores, model.Score{
			ClientID: score.ClientID,
			Score:    score.Score,
		})
	}
	leaderboard := model.Leaderboard{
		TopPlayers: scores,
	}
	err = j.RDB.SetLeaderboard(ctx, leaderboard)
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
	c.AddJob("@every 5s", ResetLeaderboard{DB, RDB})
	wg.Add(1)
	wg.Wait()
}
