package models

import (
	"time"
)

type URL struct {
	ID        string    `json:"id" db:"id"`
	LongURL   string    `json:"long_url" db:"long_url"`
	ShortCode string    `json:"short_code" db:"short_code"`
	Visits    int64     `json:"visits" db:"visits"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type URLInput struct {
	LongURL   string    `json:"long_url" binding:"required,url"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
