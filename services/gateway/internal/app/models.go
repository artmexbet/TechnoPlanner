package app

//easyjson:json
type UserRequest struct {
	Username string `json:"username" validate:"min=3,max=30"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=6"`
}
