package model

import "time"

type SessionStatus string

const (
	SessionStatusScheduled SessionStatus = "SCHEDULED"
	SessionStatusActive    SessionStatus = "ACTIVE"
	SessionStatusClosed    SessionStatus = "CLOSED"
	SessionStatusCancelled SessionStatus = "CANCELLED"
)

func IsValidSessionStatus(s string) bool {
	switch SessionStatus(s) {
	case SessionStatusScheduled, SessionStatusActive, SessionStatusClosed, SessionStatusCancelled:
		return true
	default:
		return false
	}
}

type LabSession struct {
	ID          string        `json:"id"`
	Code        string        `json:"code"`
	Title       string        `json:"title"`
	CourseCode  *string       `json:"course_code,omitempty"`
	Location    string        `json:"location"`
	SessionDate string        `json:"session_date"` // YYYY-MM-DD
	StartAt     time.Time     `json:"start_at"`     // UTC
	EndAt       time.Time     `json:"end_at"`       // UTC
	Status      SessionStatus `json:"status"`
	CreatedBy   *string       `json:"created_by,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type SessionResponse struct {
	ID          string        `json:"id"`
	Code        string        `json:"code"`
	Title       string        `json:"title"`
	CourseCode  *string       `json:"course_code,omitempty"`
	Location    string        `json:"location"`
	SessionDate string        `json:"session_date"`
	StartAt     string        `json:"start_at"` // RFC3339
	EndAt       string        `json:"end_at"`   // RFC3339
	Status      SessionStatus `json:"status"`
	CreatedBy   *string       `json:"created_by,omitempty"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type CreateSessionRequest struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	CourseCode  *string `json:"course_code,omitempty"`
	Location    string  `json:"location"`
	SessionDate string  `json:"session_date"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`
	Status      *string `json:"status,omitempty"`
}

type UpdateSessionRequest struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	CourseCode  *string `json:"course_code,omitempty"`
	Location    string  `json:"location"`
	SessionDate string  `json:"session_date"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`
	Status      *string `json:"status,omitempty"`
}

type UpdateSessionStatusRequest struct {
	Status string `json:"status"`
}

type SessionListQuery struct {
	Page        int
	Limit       int
	Status      string
	SessionDate string
	Search      string
}

type SessionListResponse struct {
	Data []SessionResponse `json:"data"`
	Meta PaginationMeta    `json:"meta"`
}
