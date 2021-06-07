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
	Describe("test SetLeaderBoard", func() {
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
})
