package models

import "time"

// Role represents the role of a user
type Role string

const (
	Admin    Role = "admin"
	SubAdmin Role = "sub-admin"
	User     Role = "user"
)

// Users User represents a user in the handlers
type Users struct {
	ID           uint      `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedBy    *uint     `db:"created_by"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Address represents an address in the handlers
type Address struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	UserID    uint      `db:"user_id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	CreatedAt time.Time `db:"created_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type UserRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddressRequestBody struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

type Session struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	SessionID string    `json:"session_id" db:"session_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Response defines the structure for API responses.
type Response struct {
	Message string      `json:"message,omitempty"` // For success responses
	Code    int         `json:"code"`              // HTTP status code
	Data    interface{} `json:"data,omitempty"`    // Actual data payload for success responses
	Error   string      `json:"error,omitempty"`   // Error message for failure responses
}
