package testUtil

import (
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateScores(count int) []repository.ScoreORM {
	var scores []repository.ScoreORM
	for i := count; i > 0; i-- {
		clientID := uuid.NewString()
		scores = append(scores, repository.ScoreORM{
			ClientID: clientID,
			Score:    float64(i),
		})
	}
	return scores
}

func CreateLeaderboard(scoreObjs []repository.ScoreORM) model.Leaderboard {
	var scores []model.Score
	for _, score := range scoreObjs {
		scores = append(scores, model.Score{
			ClientID: score.ClientID,
			Score:    score.Score,
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
