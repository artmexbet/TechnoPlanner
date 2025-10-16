package models

import (
	"proto"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password string) (bool, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), passwordHash)
	if err != nil {
		return false, nil
	}
	return true, nil
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	DeviceID  string `json:"device_id"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
}

func UserLoginFromProto(in *proto.LoginRequest) *LoginRequest {
	return &LoginRequest{
		Username:  in.GetUsername(),
		Password:  in.GetPassword(),
		DeviceID:  in.GetDeviceId(),
		UserAgent: in.GetUserAgent(),
		IP:        in.GetIpAddress(),
	}
}
