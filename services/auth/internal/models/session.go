package models

import (
	"time"
)

type Session struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	DeviceID  string    `json:"device_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
