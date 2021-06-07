package repository_test

import (
	"context"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func setupDB() *repository.Repository {
	db := repository.Init(10, 1)
	db.AutoMigrate(&model.Score{})
	return db
}

func tearDownDB(db *repository.Repository) {
	db.Migrator().DropTable(&model.Score{})
}

func truncateTable(db *repository.Repository) {
	var modelArray []interface{}
	modelArray = append(modelArray, &model.Score{})
	sql := "TRUNCATE TABLE  %s;"
	stmt := &gorm.Statement{DB: db.DB}
	for _, m := range modelArray {
		stmt.Parse(m)
		db.Exec(fmt.Sprintf(sql, stmt.Table))
	}
}

func insertScore(db *repository.Repository, score *model.Score) model.Score {
	db.Create(score)
	return *score
}

var _ = Describe("Model test", func() {
	var db *repository.Repository
	BeforeEach(func() {
		db = setupDB()
	})
	AfterEach(func() {
		tearDownDB(db)
	})
	Describe("test ListTopScores", func() {
		Context("successful", func() {
			var expectedScores []model.Score
			BeforeEach(func() {
				expectedScores = []model.Score{}
				for i := 20; i > 0; i-- {
					clientId := uuid.NewString()
					score := &model.Score{
						Score:    int64(i),
						ClientID: clientId,
					}
					expectedScores = append(expectedScores, insertScore(db, score))
				}
			})
			AfterEach(func() {
				truncateTable(db)
			})
			It("Found without limit value (default is 10)", func() {
				ctx := context.Background()
				scores, err := db.ListTopScores(ctx)
				Expect(len(scores)).To(Equal(10))
				for idx, score := range scores {
					Expect(score.ID).To(Equal(expectedScores[idx].ID))
					Expect(score.ClientID).To(Equal(expectedScores[idx].ClientID))
					Expect(score.Score).To(Equal(expectedScores[idx].Score))
				}
				Expect(err).To(BeNil())
			})
			It("Found with specific limit", func() {
				ctx := context.Background()
				for _, limit := range []int{5, 10, 20} {
					fmt.Println(limit)
					scores, err := db.ListTopScores(ctx, limit)
					Expect(len(scores)).To(Equal(limit))
					for idx, score := range scores {
						fmt.Println(score)
						fmt.Println(expectedScores[idx])
						Expect(score.ID).To(Equal(expectedScores[idx].ID))
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
				Expect(scores).To(Equal([]model.Score{}))
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})
	})
})
