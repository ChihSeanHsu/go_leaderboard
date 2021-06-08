package main_test

import (
	"context"
	"encoding/json"
	"example.com/leaderboard/cmd/web"
	"example.com/leaderboard/internal/testUtil"
	"example.com/leaderboard/pkg/cache"
	"example.com/leaderboard/pkg/model"
	"example.com/leaderboard/pkg/repository"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

func request(router *gin.Engine, method string, uri string, body io.Reader, headers ...Header) *httptest.ResponseRecorder {
	r := httptest.NewRecorder()
	req, _ := http.NewRequest(method, uri, body)
	if len(headers) > 0 {
		for _, header := range headers {
			req.Header.Add(header.Key, header.Value)
		}
	}
	router.ServeHTTP(r, req)
	return r
}

var _ = Describe("Web", func() {
	var rdb *cache.Cache
	var db *repository.Repository
	var router *gin.Engine
	BeforeEach(func() {
		rdb = cache.Init(1)
		db = testUtil.SetupDB()
		router = main.SetupRouter()
	})
	AfterEach(func() {
		testUtil.TearDownDB(db)
		rdb.Close()
	})
	Describe("test GetLeaderboard", func() {
		const (
			method = "GET"
			uri    = "/api/v1/leaderboard"
		)
		Context("successful", func() {
			var expectedLeaderboard model.Leaderboard
			BeforeEach(func() {
				scores := testUtil.CreateScores(10)
				ctx := context.Background()
				expectedLeaderboard, _ = rdb.SetLeaderboard(ctx, scores)
			})
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("Just GetLeaderboard", func() {
				r := request(router, method, uri, nil)
				expected, _ := json.Marshal(expectedLeaderboard)
				Expect(r.Code).To(Equal(http.StatusOK))
				Expect(r.Body.Bytes()).To(Equal(expected))
			})
		})
		Context("not found", func() {
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
				testUtil.TruncateTable(db)
			})
			It("no record", func() {
				r := request(router, method, uri, nil)
				expected, _ := json.Marshal(model.Leaderboard{})
				Expect(r.Code).To(Equal(http.StatusOK))
				Expect(r.Body.Bytes()).To(Equal(expected))

				scores := testUtil.CreateScores(5)
				expectedLeaderboard := testUtil.CreateLeaderboard(scores)
				expectedJSON, _ := json.Marshal(expectedLeaderboard)
				ctx := context.Background()
				for _, score := range expectedLeaderboard.TopPlayers {
					db.CreateScore(ctx, score.ClientID, score.Score)
				}
				r = request(router, method, uri, nil)
				Expect(r.Code).To(Equal(http.StatusOK))
				Expect(r.Body.Bytes()).To(Equal(expectedJSON))

			})
			It("no cache", func() {
				scores := testUtil.CreateScores(10)
				expectedLeaderboard := testUtil.CreateLeaderboard(scores)
				expectedJSON, _ := json.Marshal(expectedLeaderboard)
				ctx := context.Background()
				for _, score := range expectedLeaderboard.TopPlayers {
					db.CreateScore(ctx, score.ClientID, score.Score)
				}
				r := request(router, method, uri, nil)
				Expect(r.Code).To(Equal(http.StatusOK))
				Expect(r.Body.Bytes()).To(Equal(expectedJSON))
			})
		})
		Context("failed", func() {
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("internal error", func() {
				ctx := context.Background()
				rdb.Set(ctx, cache.LeaderboardKey, "test", 0)
				r := request(router, method, uri, nil)

				Expect(r.Code).To(Equal(http.StatusInternalServerError))
				Expect(r.Body.String()).To(Equal("{\"status\":\"Server error\"}"))
			})
		})
	})
	Describe("test PostScore", func() {
		const (
			method = "POST"
			uri    = "/api/v1/score"
		)
		Context("successful", func() {
			AfterEach(func() {
				testUtil.TruncateTable(db)
			})
			It("Just PostScore", func() {
				body := strings.NewReader("{\"score\": 10}")
				header := Header{
					Key:   "clientId",
					Value: "test",
				}
				r := request(router, method, uri, body, header)
				Expect(r.Code).To(Equal(http.StatusOK))
				Expect(r.Body.String()).To(Equal("{\"status\":\"ok\"}"))
			})
		})
		Context("failed", func() {
			AfterEach(func() {
				ctx := context.Background()
				rdb.Del(ctx, cache.LeaderboardKey)
			})
			It("without clientId", func() {
				body := strings.NewReader("{\"score\": 10}")
				r := request(router, method, uri, body)

				Expect(r.Code).To(Equal(http.StatusForbidden))
				Expect(r.Body.String()).To(Equal("{\"status\":\"No client ID\"}"))
			})
			It("invalid body", func() {
				body := strings.NewReader("{\"score\": \"test\"}")
				header := Header{
					Key:   "clientId",
					Value: "test",
				}
				r := request(router, method, uri, body, header)
				Expect(r.Code).To(Equal(http.StatusForbidden))
				Expect(r.Body.String()).To(Equal("{\"status\":\"Invalid Request Body\"}"))

				body = strings.NewReader("{\"score\": \"10\"}")
				r = request(router, method, uri, body, header)
				Expect(r.Code).To(Equal(http.StatusForbidden))
				Expect(r.Body.String()).To(Equal("{\"status\":\"Invalid Request Body\"}"))

				body = strings.NewReader("{\"test\": 10}")
				r = request(router, method, uri, body, header)
				Expect(r.Code).To(Equal(http.StatusForbidden))
				Expect(r.Body.String()).To(Equal("{\"status\":\"Invalid Request Body\"}"))

				r = request(router, method, uri, nil, header)
				Expect(r.Code).To(Equal(http.StatusForbidden))
				Expect(r.Body.String()).To(Equal("{\"status\":\"Invalid Request Body\"}"))
			})
		})
	})
})
