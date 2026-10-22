package model

import "time"

type AuditLog struct {
	ID         string                 `json:"id"`
	UserID     *string                `json:"user_id"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entity_type"`
	EntityID   string                 `json:"entity_id"`
	Status     string                 `json:"status"` // SUCCESS, FAILED
	IPAddress  *string                `json:"ip_address,omitempty"`
	UserAgent  *string                `json:"user_agent,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`

	// Relational details
	UserIdentifier string `json:"user_identifier,omitempty"`
	UserFullName   string `json:"user_full_name,omitempty"`
}

type AuditLogResponse struct {
	ID             string                 `json:"id"`
	UserID         *string                `json:"user_id,omitempty"`
	Action         string                 `json:"action"`
	EntityType     string                 `json:"entity_type"`
	EntityID       string                 `json:"entity_id"`
	Status         string                 `json:"status"`
	IPAddress      *string                `json:"ip_address,omitempty"`
	UserAgent      *string                `json:"user_agent,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"created_at"`
	UserIdentifier string                 `json:"user_identifier,omitempty"`
	UserFullName   string                 `json:"user_full_name,omitempty"`
}

type AuditLogListQuery struct {
	Page       int
	Limit      int
	Action     string
	EntityType string
	UserID     string
}

type AuditLogListResponse struct {
	Data []AuditLogResponse `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}
