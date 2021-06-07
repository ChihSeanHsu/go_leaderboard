package repository

import (
	"context"
	"example.com/leaderboard/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Repository struct {
	*gorm.DB
}

func Init(pool int, retry int) *Repository {
	dsn := os.Getenv("DB_CONN_STR")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil && retry <= 3 {
		log.Println(err)
		// waiting for return
		waitSec := 10 * retry
		retry++
		log.Println("wait for reconnect...")
		time.Sleep(time.Duration(waitSec) * time.Second)
		return Init(pool, retry)
	} else if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(pool)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	return &Repository{db}
}

func (repo *Repository) ListTop10Scores(ctx context.Context) ([]model.Score, error) {
	var scores []model.Score

	return scores, nil
}
