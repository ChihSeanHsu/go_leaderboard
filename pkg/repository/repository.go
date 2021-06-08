package repository

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repo *Repository) CreateScore(ctx context.Context, clientID string, scorePoint float64) error {
	score := &ScoreORM{
		ClientID: clientID,
		Score:    scorePoint,
	}
	sess := repo.Session(&gorm.Session{})
	err := sess.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "client_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at", "score"}),
		},
	).Create(score).Error
	return err
}

func (repo *Repository) ListTopScores(ctx context.Context, limit ...int) ([]ScoreORM, error) {
	var scores []ScoreORM
	var l int
	if len(limit) == 0 {
		l = 10
	} else {
		l = limit[0]
	}
	err := repo.Select("client_id", "score").Order("score desc").Limit(l).Find(&scores).Error
	if len(scores) == 0 {
		return scores, ErrNotFound
	}
	return scores, err
}
