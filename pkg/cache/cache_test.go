package cache_test

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createLeaderboard() model.Leaderboard {
	var scores []model.Score
	for i := 10; i > 0; i-- {
		clientID := uuid.NewString()
		scores = append(scores, model.Score{
			ClientID: clientID,
			Score:    int64(i),
		})
	}
	return model.Leaderboard{
		TopPlayers: scores,
	}
}

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
			var leaderboard model.Leaderboard
			BeforeEach(func() {
				leaderboard = createLeaderboard()
			})
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})

			It("Just set leaderboard", func() {
				var actualLeaderboard model.Leaderboard
				ctx := context.Background()
				err := rdb.SetLeaderboard(ctx, leaderboard)
				Expect(err).To(BeNil())
				actualValue, _ := rdb.Get(ctx, cache.LeaderboardKey).Result()
				json.Unmarshal([]byte(actualValue), &actualLeaderboard)
				Expect(actualLeaderboard).To(Equal(leaderboard))
			})
			It("set leaderboard twice", func() {
				var actualLeaderboard model.Leaderboard
				ctx := context.Background()
				err := rdb.SetLeaderboard(ctx, leaderboard)
				Expect(err).To(BeNil())

				leaderboard2 := createLeaderboard()
				err = rdb.SetLeaderboard(ctx, leaderboard2)
				Expect(err).To(BeNil())

				actualValue, _ := rdb.Get(ctx, cache.LeaderboardKey).Result()
				json.Unmarshal([]byte(actualValue), &actualLeaderboard)

				Expect(actualLeaderboard).NotTo(Equal(leaderboard))
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
				expectedLeaderboard = createLeaderboard()
				rdb.SetLeaderboard(ctx, expectedLeaderboard)
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
