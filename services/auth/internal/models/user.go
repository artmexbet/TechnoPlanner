package models

import (
	"time"

	"proto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	RoleID       int32     `json:"role_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
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

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
}

func UserRegisterFromProto(in *proto.RegisterRequest) *RegisterRequest {
	return &RegisterRequest{
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
}
