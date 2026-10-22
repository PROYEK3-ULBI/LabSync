package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type User struct {
	ID           string    `json:"id"`
	RoleID       string    `json:"role_id"`
	RoleCode     string    `json:"role"`
	Identifier   string    `json:"identifier"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"-"` // Tidak boleh diekspos ke JSON
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID           string     `json:"id"`
	Identifier   string     `json:"identifier"`
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	FaceEnrolled *bool      `json:"face_enrolled,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type CreateUserRequest struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Status     string `json:"status"`
}

type UpdateUserRequest struct {
	Identifier string  `json:"identifier"`
	Email      string  `json:"email"`
	FullName   string  `json:"full_name"`
	Password   *string `json:"password,omitempty"`
	Role       *string `json:"role,omitempty"`
	Status     *string `json:"status,omitempty"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status"`
}

type UserListQuery struct {
	Page   int
	Limit  int
	Search string
	Role   string
	Status string
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type UserListResponse struct {
	Data []UserResponse `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type Claims struct {
	UserID     string `json:"user_id"`
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}
