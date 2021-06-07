package main

import (
	"example.com/leaderboard/pkg/repository"
	"github.com/gin-gonic/gin"
)

var (
	DB *repository.Repository
)

func main() {
	DB = repository.Init(20, 1)
	r := gin.Default()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
