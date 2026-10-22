package model

import "time"

type Attendance struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	SessionID          string     `json:"session_id"`
	CheckedInAt        time.Time  `json:"checked_in_at"`
	CheckedOutAt       *time.Time `json:"checked_out_at,omitempty"`
	VerificationStatus string     `json:"verification_status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Informasi relasi untuk pembacaan list/history
	UserIdentifier string `json:"user_identifier,omitempty"`
	UserFullName   string `json:"user_full_name,omitempty"`
	SessionCode    string `json:"session_code,omitempty"`
	SessionTitle   string `json:"session_title,omitempty"`
	SessionDate    string `json:"session_date,omitempty"`
}

type AttendanceResponse struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"user_id"`
	SessionID          string  `json:"session_id"`
	CheckedInAt        string  `json:"checked_in_at"`
	CheckedOutAt       *string `json:"checked_out_at,omitempty"`
	VerificationStatus string  `json:"verification_status"`
	UserIdentifier     string  `json:"user_identifier,omitempty"`
	UserFullName       string  `json:"user_full_name,omitempty"`
	SessionCode        string  `json:"session_code,omitempty"`
	SessionTitle       string  `json:"session_title,omitempty"`
	SessionDate        string  `json:"session_date,omitempty"`
}

type CheckInResponse struct {
	Message    string             `json:"message"`
	Attendance AttendanceResponse `json:"attendance"`
}

type CheckOutRequest struct {
	SessionID string `json:"session_id"`
}

type CheckOutResponse struct {
	Message    string             `json:"message"`
	Attendance AttendanceResponse `json:"attendance"`
}

type AttendanceListQuery struct {
	Page      int
	Limit     int
	SessionID string
	UserID    string
	Date      string
}

type AttendanceListResponse struct {
	Data []AttendanceResponse `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}
