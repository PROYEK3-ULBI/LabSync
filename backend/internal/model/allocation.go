package model

import "time"

type Allocation struct {
	ID              string     `json:"id"`
	AssetID         string     `json:"asset_id"`
	UserID          string     `json:"user_id"`
	SessionID       string     `json:"session_id"`
	AllocatedAt     time.Time  `json:"allocated_at"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty"`
	ConditionBefore string     `json:"condition_before"`
	ConditionAfter  *string    `json:"condition_after,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relational details
	AssetCode      string `json:"asset_code,omitempty"`
	AssetName      string `json:"asset_name,omitempty"`
	UserIdentifier string `json:"user_identifier,omitempty"`
	UserFullName   string `json:"user_full_name,omitempty"`
	SessionCode    string `json:"session_code,omitempty"`
	SessionTitle   string `json:"session_title,omitempty"`
}

type AllocationResponse struct {
	ID              string  `json:"id"`
	AssetID         string  `json:"asset_id"`
	UserID          string  `json:"user_id"`
	SessionID       string  `json:"session_id"`
	AllocatedAt     string  `json:"allocated_at"`
	ReturnedAt      *string `json:"returned_at,omitempty"`
	ConditionBefore string  `json:"condition_before"`
	ConditionAfter  *string `json:"condition_after,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	AssetCode       string  `json:"asset_code,omitempty"`
	AssetName       string  `json:"asset_name,omitempty"`
	UserIdentifier  string  `json:"user_identifier,omitempty"`
	UserFullName    string  `json:"user_full_name,omitempty"`
	SessionCode     string  `json:"session_code,omitempty"`
	SessionTitle    string  `json:"session_title,omitempty"`
}

type CreateAllocationRequest struct {
	AssetID         string  `json:"asset_id"`
	UserID          string  `json:"user_id"`
	SessionID       string  `json:"session_id"`
	ConditionBefore *string `json:"condition_before,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

type ReturnAllocationRequest struct {
	Condition string  `json:"condition"` // GOOD | DAMAGED
	Notes     *string `json:"notes,omitempty"`
}

type AllocationListQuery struct {
	Page      int
	Limit     int
	AssetID   string
	UserID    string
	SessionID string
	Status    string // ACTIVE (belum dikembalikan), RETURNED (sudah dikembalikan)
}

type AllocationListResponse struct {
	Data []AllocationResponse `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}

type AssetTimelineEvent struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"event_type"` // ALLOCATION, RETURN, MAINTENANCE, STATUS_CHANGE
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   string                 `json:"timestamp"`
	ActorID     *string                `json:"actor_id,omitempty"`
	ActorName   *string                `json:"actor_name,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type AssetHistoryResponse struct {
	Asset    AssetResponse        `json:"asset"`
	Timeline []AssetTimelineEvent `json:"timeline"`
}
