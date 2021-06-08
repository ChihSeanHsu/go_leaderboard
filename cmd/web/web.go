package main

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/internal/logging"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

var (
	DB     *repository.Repository
	Cache  *cache.Cache
	logger = logging.ZapL
)

type ScoreRequest struct {
	Score float64
}

func traceContext(c *gin.Context) context.Context {
	traceID := c.Request.Header.Get("traceId")
	if traceID == "" {
		traceID = uuid.New().String()
		c.Request.Header.Add("X-TRACE-ID", traceID)
	}
	return context.WithValue(context.Background(), "TraceID", traceID)
}

func GetLeaderboard(c *gin.Context) {
	var (
		status int
		result interface{}
		topTen []repository.ScoreORM
	)
	ctx := traceContext(c)

	leaderboard, err := Cache.GetLeaderboard(ctx)
	if err == cache.ErrNotFound || (err == nil && leaderboard.TopPlayers == nil) {
		topTen, err = DB.ListTopScores(ctx)
		if err == nil || err == repository.ErrNotFound {
			leaderboard, err = Cache.SetLeaderboard(ctx, topTen)
		}
	}
	if err != nil {
		logger(ctx).Error(err)
		status = http.StatusInternalServerError
		result = gin.H{
			"status": "Server error",
		}
	} else {
		status = http.StatusOK
		result = leaderboard
	}
	c.JSON(status, result)
}

func PostScore(c *gin.Context) {
	var (
		status    int
		result    interface{}
		scoreJSON ScoreRequest
		err       error
		body      []byte
	)
	ctx := traceContext(c)

	clientID := c.Request.Header.Get("clientId")
	if c.Request.Body != nil {
		body, err = ioutil.ReadAll(c.Request.Body)
		err = json.Unmarshal(body, &scoreJSON)
	}
	switch {
	case clientID == "":
		logger(ctx).Warn(err)
		status = http.StatusForbidden
		result = gin.H{
			"status": "No client ID",
		}
	case err != nil || scoreJSON == ScoreRequest{}:
		logger(ctx).Warn(err)
		status = http.StatusForbidden
		result = gin.H{
			"status": "Invalid Request Body",
		}
	default:
		err = DB.CreateScore(ctx, clientID, scoreJSON.Score)
		if err != nil {
			logger(ctx).Error(err)
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
	logging.InitLogging()
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
