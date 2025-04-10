package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritikchawla/url-shortner/db"
	"github.com/ritikchawla/url-shortner/models"
)

const (
	shortCodeLength = 6
	cacheDuration   = 24 * time.Hour
)

// GenerateShortCode generates a random short code
func generateShortCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)[:shortCodeLength], nil
}

// CreateShortURL handles the creation of short URLs
func CreateShortURL(c *gin.Context) {
	var input models.URLInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate short code
	shortCode, err := generateShortCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short code"})
		return
	}

	// Set default expiration if not provided
	if input.ExpiresAt.IsZero() {
		input.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // Default 7 days
	}

	// Insert into database
	query := `
INSERT INTO urls (long_url, short_code, expires_at)
VALUES ($1, $2, $3)
RETURNING id, created_at`

	var url models.URL
	err = db.DB.QueryRow(query, input.LongURL, shortCode, input.ExpiresAt).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	// Cache the URL
	err = db.CacheURL(shortCode, input.LongURL, cacheDuration)
	if err != nil {
		// Log the error but don't fail the request
	}

	url.LongURL = input.LongURL
	url.ShortCode = shortCode
	url.ExpiresAt = input.ExpiresAt

	c.JSON(http.StatusCreated, url)
}

// RedirectURL handles the redirection of short URLs
func RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Try to get URL from cache first
	longURL, err := db.GetCachedURL(shortCode)
	if err == nil && longURL != "" {
		// Increment visits asynchronously
		go func() {
			db.IncrementVisits(shortCode)
		}()
		c.Redirect(http.StatusMovedPermanently, longURL)
		return
	}

	// If not in cache, get from database
	query := `
SELECT long_url, expires_at
FROM urls
WHERE short_code = $1`

	var url models.URL
	err = db.DB.QueryRow(query, shortCode).Scan(&url.LongURL, &url.ExpiresAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}

	// Check if URL has expired
	if !url.ExpiresAt.IsZero() && url.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		return
	}

	// Cache the URL for future requests
	db.CacheURL(shortCode, url.LongURL, cacheDuration)

	// Increment visits asynchronously
	go func() {
		db.IncrementVisits(shortCode)
	}()

	c.Redirect(http.StatusMovedPermanently, url.LongURL)
}

// GetURLStats returns statistics for a short URL
func GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var url models.URL
	query := `
SELECT id, long_url, short_code, visits, created_at, expires_at
FROM urls
WHERE short_code = $1`

	err := db.DB.QueryRow(query, shortCode).Scan(
		&url.ID, &url.LongURL, &url.ShortCode,
		&url.Visits, &url.CreatedAt, &url.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL stats"})
		return
	}

	// Get real-time visit count from Redis
	visits, err := db.GetVisits(shortCode)
	if err == nil && visits > 0 {
		url.Visits = visits
	}

	c.JSON(http.StatusOK, url)
}
