package models

//go:generate easyjson -all models.go

type TokenState string

const (
	TokenStateValid   TokenState = "VALID"
	TokenStateExpired TokenState = "EXPIRED"
	TokenStateInvalid TokenState = "INVALID"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"min=3,max=30"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type LoginRequest struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	DeviceID  string `json:"-" validate:"-"`
	IP        string `json:"-" validate:"-"`
	UserAgent string `json:"-" validate:"-"`
}

//easyjson:json
type TokenPair struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenValidationResponse struct {
	UserID string     `json:"user_id"`
	State  TokenState `json:"state"`
	Role   string     `json:"role"`
}

type TokenRefreshRequest struct {
	TokenPair
	DeviceID  string `json:"-" validate:"-"`
	UserAgent string `json:"-" validate:"-"`
	IP        string `json:"-" validate:"-"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type EquipmentCategory struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedBy   *string `json:"created_by,omitempty"`
	UpdatedBy   *string `json:"updated_by,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

type Equipment struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name" validate:"required"`
	Description *string             `json:"description,omitempty"`
	Quantity    int32               `json:"quantity" validate:"gte=0"`
	Categories  []EquipmentCategory `json:"categories,omitempty"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	CreatedBy   *string             `json:"created_by,omitempty"`
	UpdatedBy   *string             `json:"updated_by,omitempty"`
	DeletedAt   *string             `json:"deleted_at,omitempty"`
}

type EquipmentCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Quantity    int32   `json:"quantity" validate:"gte=0"`
	CategoryIDs []int   `json:"category_ids" validate:"dive,gt=0"`
}

type EquipmentUpdateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Quantity    int32   `json:"quantity" validate:"gte=0"`
	CategoryIDs []int   `json:"category_ids" validate:"dive,gt=0"`
}

type EquipmentCategoryCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type EquipmentCategoryUpdateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type EquipmentCategoryListResponse struct {
	Items []EquipmentCategory `json:"items"`
}

type RequestResponse struct {
	ID            string      `json:"id"`
	RequestText   *string     `json:"request_text,omitempty"`
	Status        string      `json:"status"`
	ScheduleTime  string      `json:"schedule_time"`
	EndTime       string      `json:"end_time"`
	Address       string      `json:"address"`
	ResponsibleID *string     `json:"responsible_id,omitempty"`
	Equipment     []Equipment `json:"equipment,omitempty"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	CreatedBy     *string     `json:"created_by,omitempty"`
	UpdatedBy     *string     `json:"updated_by,omitempty"`
	DeletedAt     *string     `json:"deleted_at,omitempty"`
}

type RequestListResponse struct {
	Items []RequestResponse `json:"items"`
}

type RequestFilter struct {
	ResponsibleID *string `json:"responsible_id,omitempty"`
}

type RequestStatusHistoryResponse struct {
	ID        int32   `json:"id"`
	RequestID string  `json:"request_id"`
	Status    string  `json:"status"`
	Comment   *string `json:"comment,omitempty"`
	ChangedBy *string `json:"changed_by,omitempty"`
	ChangedAt string  `json:"changed_at"`
}

type RequestStatusHistoryListResponse struct {
	Items []RequestStatusHistoryResponse `json:"items"`
}

type RequestStatusUpdateRequest struct {
	Status  string  `json:"status" validate:"required,oneof=canceled pending assigned in_progress completed rejected"`
	Comment *string `json:"comment,omitempty"`
}

type PorterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PorterListResponse struct {
	Items []PorterResponse `json:"items"`
}

type PorterCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type MeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
