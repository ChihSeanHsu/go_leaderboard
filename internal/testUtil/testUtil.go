package testUtil

import (
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateLeaderboard() model.Leaderboard {
	var scores []model.Score
	for i := 10; i > 0; i-- {
		clientID := uuid.NewString()
		scores = append(scores, model.Score{
			ClientID: clientID,
			Score:    float64(i),
		})
	}
	return model.Leaderboard{
		TopPlayers: scores,
	}
}

func SetupDB() *repository.Repository {
	db := repository.Init(10, 1)
	db.AutoMigrate(&repository.ScoreORM{})
	return db
}

func TearDownDB(db *repository.Repository) {
	db.Migrator().DropTable(&repository.ScoreORM{})
}

func TruncateTable(db *repository.Repository) {
	var modelArray []interface{}
	modelArray = append(modelArray, &repository.ScoreORM{})
	sql := "TRUNCATE TABLE  %s;"
	stmt := &gorm.Statement{DB: db.DB}
	for _, m := range modelArray {
		stmt.Parse(m)
		db.Exec(fmt.Sprintf(sql, stmt.Table))
	}
}
