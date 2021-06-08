package main

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/repository"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	DB    *repository.Repository
	Cache *cache.Cache
)

type ScoreRequest struct {
	Score int64
}

func GetLeaderboard(c *gin.Context) {
	var status int
	var result interface{}
	ctx := context.Background()
	leaderboard, err := Cache.GetLeaderboard(ctx)
	switch {
	case err == cache.ErrNotFound:
		status = http.StatusNotFound
		result = gin.H{
			"status": "leaderboard not ready",
		}
	case err != nil:
		status = http.StatusInternalServerError
		result = gin.H{
			"status": "Server error",
		}
	default:
		status = http.StatusOK
		result = leaderboard
	}
	c.JSON(status, result)
}

func PostScore(c *gin.Context) {
	var status int
	var result interface{}
	var scoreJSON ScoreRequest
	var err error
	var body []byte
	clientID := c.Request.Header.Get("clientId")
	if c.Request.Body != nil {
		body, err = ioutil.ReadAll(c.Request.Body)
		err = json.Unmarshal(body, &scoreJSON)
	}
	switch {
	case clientID == "":
		status = http.StatusForbidden
		result = gin.H{
			"status": "No client ID",
		}
	case err != nil || scoreJSON == ScoreRequest{}:
		log.Println(err)
		status = http.StatusForbidden
		result = gin.H{
			"status": "Invalid Request Body",
		}
	default:
		ctx := context.Background()
		err = DB.CreateScore(ctx, clientID, scoreJSON.Score)
		if err != nil {
			log.Println(err)
			status = http.StatusInternalServerError
			result = gin.H{
				"status": "Store Score error",
			}
		} else {
			status = http.StatusOK
			result = gin.H{
				"status": "ok",
			}
		}
	}
	c.JSON(status, result)
}

func SetupRouter() *gin.Engine {
	DB = repository.Init(20, 1)
	Cache = cache.Init(20)
	r := gin.Default()
	r.GET("/api/v1/leaderboard", GetLeaderboard)
	r.POST("/api/v1/score", PostScore)
	return r
}

func main() {
	r := SetupRouter()
	defer Cache.Close()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
