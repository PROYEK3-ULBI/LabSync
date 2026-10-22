package model

import "time"

type MaintenanceStatus string

const (
	MaintenanceStatusOpen       MaintenanceStatus = "OPEN"
	MaintenanceStatusInProgress MaintenanceStatus = "IN_PROGRESS"
	MaintenanceStatusCompleted  MaintenanceStatus = "COMPLETED"
	MaintenanceStatusCancelled  MaintenanceStatus = "CANCELLED"
)

func IsValidMaintenanceStatus(s string) bool {
	switch MaintenanceStatus(s) {
	case MaintenanceStatusOpen, MaintenanceStatusInProgress, MaintenanceStatusCompleted, MaintenanceStatusCancelled:
		return true
	default:
		return false
	}
}

type MaintenancePriority string

const (
	MaintenancePriorityLow    MaintenancePriority = "LOW"
	MaintenancePriorityMedium MaintenancePriority = "MEDIUM"
	MaintenancePriorityHigh   MaintenancePriority = "HIGH"
)

func IsValidMaintenancePriority(p string) bool {
	switch MaintenancePriority(p) {
	case MaintenancePriorityLow, MaintenancePriorityMedium, MaintenancePriorityHigh:
		return true
	default:
		return false
	}
}

type Maintenance struct {
	ID          string              `json:"id"`
	AssetID     string              `json:"asset_id"`
	ReportedBy  string              `json:"reported_by"`
	IssueType   string              `json:"issue_type"`
	Description string              `json:"description"`
	Priority    MaintenancePriority `json:"priority"`
	Status      MaintenanceStatus   `json:"status"`
	StartedAt   *time.Time          `json:"started_at,omitempty"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	Resolution  *string             `json:"resolution,omitempty"`
	Cost        *float64            `json:"cost,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`

	// Relational details
	AssetCode          string `json:"asset_code,omitempty"`
	AssetName          string `json:"asset_name,omitempty"`
	ReporterIdentifier string `json:"reporter_identifier,omitempty"`
	ReporterName       string `json:"reporter_name,omitempty"`
}

type MaintenanceResponse struct {
	ID                 string              `json:"id"`
	AssetID            string              `json:"asset_id"`
	ReportedBy         string              `json:"reported_by"`
	IssueType          string              `json:"issue_type"`
	Description        string              `json:"description"`
	Priority           MaintenancePriority `json:"priority"`
	Status             MaintenanceStatus   `json:"status"`
	StartedAt          *string             `json:"started_at,omitempty"`
	CompletedAt        *string             `json:"completed_at,omitempty"`
	Resolution         *string             `json:"resolution,omitempty"`
	Cost               *float64            `json:"cost,omitempty"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
	AssetCode          string              `json:"asset_code,omitempty"`
	AssetName          string              `json:"asset_name,omitempty"`
	ReporterIdentifier string              `json:"reporter_identifier,omitempty"`
	ReporterName       string              `json:"reporter_name,omitempty"`
}

type CreateMaintenanceRequest struct {
	AssetID     string  `json:"asset_id"`
	IssueType   *string `json:"issue_type,omitempty"`
	Description string  `json:"description"`
	Priority    *string `json:"priority,omitempty"`
}

type UpdateMaintenanceRequest struct {
	Status       *string  `json:"status,omitempty"`
	Priority     *string  `json:"priority,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Resolution   *string  `json:"resolution,omitempty"`
	Cost         *float64 `json:"cost,omitempty"`
	ResultStatus *string  `json:"result_status,omitempty"` // AVAILABLE atau RETIRED saat menyelesaikan pemeliharaan
}

type MaintenanceListQuery struct {
	Page     int
	Limit    int
	AssetID  string
	Status   string
	Priority string
}

type MaintenanceListResponse struct {
	Data []MaintenanceResponse `json:"data"`
	Meta PaginationMeta        `json:"meta"`
}
