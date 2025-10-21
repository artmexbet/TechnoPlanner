package models

//go:generate easyjson -all models.go

type TokenState string

const (
	TokenStateValid   TokenState = "valid"
	TokenStateExpired TokenState = "expired"
	TokenStateInvalid TokenState = "invalid"
)

//easyjson:json
type RegisterRequest struct {
	Username string `json:"username" validate:"min=3,max=30"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

//easyjson:json
type LoginRequest struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	DeviceID  string `json:"-" validate:"-"`
	IP        string `json:"-" validate:"-"`
	UserAgent string `json:"-" validate:"-"`
}

//easyjson:json
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenValidationResponse struct {
	UserID string     `json:"user_id"`
	State  TokenState `json:"state"`
}

type TokenRefreshRequest struct {
	Pair      TokenPair `json:"pair"`
	DeviceID  string    `json:"-" validate:"-"`
	UserAgent string    `json:"-" validate:"-"`
	IP        string    `json:"-" validate:"-"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
