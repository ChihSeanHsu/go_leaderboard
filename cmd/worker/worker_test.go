package main_test

import (
	"context"
	"example.com/leaderboard/cmd/worker"
	"example.com/leaderboard/internal/testUtil"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	var rdb *cache.Cache
	var db *repository.Repository
	BeforeEach(func() {
		rdb = cache.Init(1)
		db = testUtil.SetupDB()
	})
	AfterEach(func() {
		testUtil.TearDownDB(db)
		rdb.Close()
	})
	Describe("test ResetLeaderboard", func() {
		Context("successful", func() {
			var expectedLeaderboard model.Leaderboard
			BeforeEach(func() {
				expectedLeaderboard = testUtil.CreateLeaderboard()
				ctx := context.Background()
				for _, score := range expectedLeaderboard.TopPlayers {
					db.CreateScore(ctx, score.ClientID, score.Score)
				}
			})
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("Just reset leaderboard", func() {
				job := main.ResetLeaderboard{DB: db, RDB: rdb}
				job.Run()
				ctx := context.Background()
				leaderboard, _ := rdb.GetLeaderboard(ctx)
				Expect(leaderboard).To(Equal(expectedLeaderboard))
			})
		})
		Context("failed", func() {
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("", func() {
				job := main.ResetLeaderboard{DB: db, RDB: rdb}
				job.Run()
				ctx := context.Background()
				leaderboard, err := rdb.GetLeaderboard(ctx)
				Expect(err).To(Equal(cache.ErrNotFound))
				Expect(leaderboard).To(Equal(model.Leaderboard{}))
			})
		})
	})
})
