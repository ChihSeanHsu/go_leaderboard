package repository_test

import (
	"context"
	"example.com/leaderboard/internal/testUtil"
	"example.com/leaderboard/pkg/repository"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


func insertScore(db *repository.Repository, score *repository.ScoreORM) repository.ScoreORM {
	db.Create(score)
	return *score
}

var _ = Describe("Model test", func() {
	var db *repository.Repository
	BeforeEach(func() {
		db = testUtil.SetupDB()
	})
	AfterEach(func() {
		testUtil.TearDownDB(db)
	})
	Describe("test ListTopScores", func() {
		Context("successful", func() {
			var expectedScores []repository.ScoreORM
			BeforeEach(func() {
				expectedScores = []repository.ScoreORM{}
				for i := 20; i > 0; i-- {
					clientId := uuid.NewString()
					score := &repository.ScoreORM{
						Score:    int64(i),
						ClientID: clientId,
					}
					expectedScores = append(expectedScores, insertScore(db, score))
				}
			})
			AfterEach(func() {
				testUtil.TruncateTable(db)
			})
			It("Found without limit value (default is 10)", func() {
				ctx := context.Background()
				scores, err := db.ListTopScores(ctx)
				Expect(len(scores)).To(Equal(10))
				for idx, score := range scores {
					Expect(score.ClientID).To(Equal(expectedScores[idx].ClientID))
					Expect(score.Score).To(Equal(expectedScores[idx].Score))
				}
				Expect(err).To(BeNil())
			})
			It("Found with specific limit", func() {
				ctx := context.Background()
				for _, limit := range []int{5, 10, 20} {
					scores, err := db.ListTopScores(ctx, limit)
					Expect(len(scores)).To(Equal(limit))
					for idx, score := range scores {
						Expect(score.ClientID).To(Equal(expectedScores[idx].ClientID))
						Expect(score.Score).To(Equal(expectedScores[idx].Score))
					}
					Expect(err).To(BeNil())
				}
			})
		})
		Context("failed", func() {
			It("not found", func() {
				ctx := context.Background()
				scores, err := db.ListTopScores(ctx, 10)
				Expect(scores).To(Equal([]repository.ScoreORM{}))
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})
	})
	Describe("test CreateScore", func() {
		Context("successful", func() {
			AfterEach(func() {
				testUtil.TruncateTable(db)
			})
			It("default", func() {
				var actualScore repository.ScoreORM
				ctx := context.Background()
				clientId := uuid.NewString()
				score := int64(10)
				err := db.CreateScore(ctx, clientId, score)
				db.Where("client_id=? AND score=?", clientId, score).First(&actualScore)
				Expect(err).To(BeNil())
				Expect(actualScore.ClientID).To(Equal(clientId))
				Expect(actualScore.Score).To(Equal(score))
			})
			It("duplicated ClientID", func() {
				var actualScore repository.ScoreORM
				ctx := context.Background()
				clientId := uuid.NewString()
				score := int64(10)
				err := db.CreateScore(ctx, clientId, score)
				Expect(err).To(BeNil())

				newScore := score + 10
				err = db.CreateScore(ctx, clientId, newScore)
				Expect(err).To(BeNil())

				err = db.Where("client_id=? AND score=?", clientId, score).First(&actualScore).Error
				Expect(err).ToNot(BeNil())

				db.Where("client_id=? AND score=?", clientId, newScore).First(&actualScore)
				// check update
				Expect(actualScore.ClientID).To(Equal(clientId))
				Expect(actualScore.Score).To(Equal(newScore))
			})
		})
	})
})
