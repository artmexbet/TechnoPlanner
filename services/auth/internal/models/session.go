package models

import (
	"fmt"
	"time"
)

var (
	ErrInvalidDevice    = fmt.Errorf("invalid device")
	ErrInvalidUserAgent = fmt.Errorf("invalid user agent")
	ErrInvalidIP        = fmt.Errorf("invalid IP address")
)

type Session struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	DeviceID  string    `json:"device_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *Session) Validate(deviceID, userAgent, ip string) error {
	if s.DeviceID != deviceID {
		return ErrInvalidDevice
	}
	if s.UserAgent != userAgent {
		return ErrInvalidUserAgent
	}
	if s.IP != ip {
		return ErrInvalidIP
	}
	return nil
}
