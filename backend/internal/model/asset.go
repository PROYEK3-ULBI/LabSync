package model

import "time"

type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "AVAILABLE"
	AssetStatusAllocated   AssetStatus = "ALLOCATED"
	AssetStatusMaintenance AssetStatus = "MAINTENANCE"
	AssetStatusRetired     AssetStatus = "RETIRED"
)

func IsValidAssetStatus(s string) bool {
	switch AssetStatus(s) {
	case AssetStatusAvailable, AssetStatusAllocated, AssetStatusMaintenance, AssetStatusRetired:
		return true
	default:
		return false
	}
}

type Asset struct {
	ID           string      `json:"id"`
	AssetCode    string      `json:"asset_code"`
	AssetName    string      `json:"asset_name"`
	Category     string      `json:"category"`
	Brand        *string     `json:"brand,omitempty"`
	Model        *string     `json:"model,omitempty"`
	SerialNumber *string     `json:"serial_number,omitempty"`
	Location     string      `json:"location"`
	Status       AssetStatus `json:"status"`
	Notes        *string     `json:"notes,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type AssetResponse struct {
	ID           string      `json:"id"`
	AssetCode    string      `json:"asset_code"`
	AssetName    string      `json:"asset_name"`
	Category     string      `json:"category"`
	Brand        *string     `json:"brand,omitempty"`
	Model        *string     `json:"model,omitempty"`
	SerialNumber *string     `json:"serial_number,omitempty"`
	Location     string      `json:"location"`
	Status       AssetStatus `json:"status"`
	Notes        *string     `json:"notes,omitempty"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
}

type CreateAssetRequest struct {
	AssetCode    string  `json:"asset_code"`
	AssetName    string  `json:"asset_name"`
	Category     *string `json:"category,omitempty"`
	Brand        *string `json:"brand,omitempty"`
	Model        *string `json:"model,omitempty"`
	SerialNumber *string `json:"serial_number,omitempty"`
	Location     *string `json:"location,omitempty"`
	Status       *string `json:"status,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

type UpdateAssetRequest struct {
	AssetCode    string  `json:"asset_code"`
	AssetName    string  `json:"asset_name"`
	Category     *string `json:"category,omitempty"`
	Brand        *string `json:"brand,omitempty"`
	Model        *string `json:"model,omitempty"`
	SerialNumber *string `json:"serial_number,omitempty"`
	Location     *string `json:"location,omitempty"`
	Status       *string `json:"status,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

type UpdateAssetStatusRequest struct {
	Status string `json:"status"`
}

type AssetListQuery struct {
	Page     int
	Limit    int
	Status   string
	Category string
	Search   string
}

type AssetListResponse struct {
	Data []AssetResponse `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}
