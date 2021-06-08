package cache_test

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/internal/testUtil"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var rdb *cache.Cache
	BeforeEach(func() {
		rdb = cache.Init(1)
	})
	AfterEach(func() {
		rdb.Close()
	})
	Describe("test SetLeaderboard", func() {
		Context("successful", func() {
			var scoreObjs []repository.ScoreORM
			BeforeEach(func() {
				scoreObjs = testUtil.CreateScores(10)
			})
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})

			It("Just set leaderboard", func() {
				var actualLeaderboard model.Leaderboard
				ctx := context.Background()
				expectedLeaderboard, err := rdb.SetLeaderboard(ctx, scoreObjs)
				Expect(err).To(BeNil())
				actualValue, _ := rdb.Get(ctx, cache.LeaderboardKey).Result()
				json.Unmarshal([]byte(actualValue), &actualLeaderboard)
				Expect(actualLeaderboard).To(Equal(expectedLeaderboard))
			})
			It("set leaderboard twice", func() {
				var actualLeaderboard model.Leaderboard
				ctx := context.Background()
				expectedLeaderboard, err := rdb.SetLeaderboard(ctx, scoreObjs)
				Expect(err).To(BeNil())

				scoreObjs2 := testUtil.CreateScores(10)
				leaderboard2, err := rdb.SetLeaderboard(ctx, scoreObjs2)
				Expect(err).To(BeNil())

				actualValue, _ := rdb.Get(ctx, cache.LeaderboardKey).Result()
				json.Unmarshal([]byte(actualValue), &actualLeaderboard)

				Expect(actualLeaderboard).NotTo(Equal(expectedLeaderboard))
				Expect(actualLeaderboard).To(Equal(leaderboard2))
			})
		})
	})
	Describe("test GetLeaderboard", func() {
		Context("successful", func() {
			var actualLeaderboard model.Leaderboard
			var expectedLeaderboard model.Leaderboard
			BeforeEach(func() {
				ctx := context.Background()
				scoreObjs := testUtil.CreateScores(10)
				expectedLeaderboard, _ = rdb.SetLeaderboard(ctx, scoreObjs)
			})
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})

			It("Just get leaderboard", func() {
				var err error
				ctx := context.Background()
				actualLeaderboard, err = rdb.GetLeaderboard(ctx)
				Expect(err).To(BeNil())
				Expect(actualLeaderboard).To(Equal(expectedLeaderboard))
			})
		})
		Context("failed", func() {
			var leaderboard model.Leaderboard
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("not found", func() {
				var err error
				ctx := context.Background()
				leaderboard, err = rdb.GetLeaderboard(ctx)
				Expect(err).To(Equal(cache.ErrNotFound))
				Expect(leaderboard).To(Equal(model.Leaderboard{}))
			})
			It("empty value", func() {
				var err error
				ctx := context.Background()
				leaderboard, err = rdb.GetLeaderboard(ctx)
				Expect(err).To(Equal(cache.ErrNotFound))
				Expect(leaderboard).To(Equal(model.Leaderboard{}))
			})
			It("invalid JSON", func() {
				var err error
				ctx := context.Background()
				_ = rdb.Set(ctx, cache.LeaderboardKey, "test", 0)
				leaderboard, err = rdb.GetLeaderboard(ctx)
				Expect(err).To(Equal(cache.ErrDataCorruption))
				Expect(leaderboard).To(Equal(model.Leaderboard{}))
			})
		})
	})
})
