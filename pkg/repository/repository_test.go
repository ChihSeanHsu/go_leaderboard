package repository_test

import (
	"context"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func setupDB() *repository.Repository {
	db := repository.Init(10, 1)
	db.AutoMigrate(&model.ScoreModel{})
	return db
}

func tearDownDB(db *repository.Repository) {
	db.Migrator().DropTable(&model.ScoreModel{})
}

func truncateTable(db *repository.Repository) {
	var modelArray []interface{}
	modelArray = append(modelArray, &model.ScoreModel{})
	sql := "TRUNCATE TABLE  %s;"
	stmt := &gorm.Statement{DB: db.DB}
	for _, m := range modelArray {
		stmt.Parse(m)
		db.Raw(fmt.Sprintf(sql, stmt.Table))
	}
}

var _ = Describe("Model test", func() {
	var db *repository.Repository
	BeforeEach(func() {
		db = setupDB()
	})
	AfterEach(func() {
		tearDownDB(db)
	})
	Describe("test ", func() {
		AfterEach(func() {
			truncateTable(db)
		})
		Context("failed", func() {
			It("not found", func() {
				ctx := context.Background()
				scores, err := db.ListTop10Scores(ctx)
				Expect(scores).To(Equal([]model.ScoreModel{}))
				Expect(err).To(Equal(repository.ErrNotFound))
			})
		})
	})
})
