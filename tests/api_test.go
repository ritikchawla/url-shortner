package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritikchawla/url-shortner/api"
	"github.com/ritikchawla/url-shortner/models"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/shorten", api.CreateShortURL)
	r.GET("/api/stats/:shortCode", api.GetURLStats)
	r.GET("/:shortCode", api.RedirectURL)
	return r
}

func TestCreateShortURL(t *testing.T) {
	router := setupRouter()

	// Test valid URL
	input := models.URLInput{
		LongURL:   "https://example.com",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	inputJSON, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.URL
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ShortCode)
	assert.Equal(t, input.LongURL, response.LongURL)

	// Test invalid URL
	input.LongURL = "not-a-url"
	inputJSON, _ = json.Marshal(input)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirectURL(t *testing.T) {
	router := setupRouter()

	// Test non-existent short code
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetURLStats(t *testing.T) {
	router := setupRouter()

	// Test non-existent short code stats
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/stats/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
