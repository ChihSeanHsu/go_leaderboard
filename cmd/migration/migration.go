package main

import (
	"example.com/leaderboard/pkg/repository"
)

func main() {
	db := repository.Init(1, 1)
	db.AutoMigrate(&repository.ScoreORM{})
}
