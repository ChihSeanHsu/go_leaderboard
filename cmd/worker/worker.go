package main

import (
	"example.com/leaderboard/pkg/repository"
	"sync"
)

var (
	DB *repository.Repository
)

func main() {
	var wg sync.WaitGroup
	DB = repository.Init(20, 1)
	wg.Wait()
}
