package model

type AttendanceReportSummary struct {
	TotalAttendance int `json:"total_attendance"`
	UniqueUsers     int `json:"unique_users"`
	UniqueSessions  int `json:"unique_sessions"`
}

type AttendanceReportResponse struct {
	Summary AttendanceReportSummary `json:"summary"`
	Data    []AttendanceResponse    `json:"data"`
}

type AssetReportSummary struct {
	TotalAssets int            `json:"total_assets"`
	ByStatus    map[string]int `json:"by_status"`
	ByCategory  map[string]int `json:"by_category"`
}

type AssetReportResponse struct {
	Summary AssetReportSummary `json:"summary"`
	Data    []AssetResponse    `json:"data"`
}

type AllocationReportSummary struct {
	TotalAllocations int `json:"total_allocations"`
	ActiveCount      int `json:"active_count"`
	ReturnedCount    int `json:"returned_count"`
	DamagedCount     int `json:"damaged_count"`
}

type AllocationReportResponse struct {
	Summary AllocationReportSummary `json:"summary"`
	Data    []AllocationResponse    `json:"data"`
}

type MaintenanceReportSummary struct {
	TotalTickets    int            `json:"total_tickets"`
	OpenCount       int            `json:"open_count"`
	InProgressCount int            `json:"in_progress_count"`
	CompletedCount  int            `json:"completed_count"`
	CancelledCount  int            `json:"cancelled_count"`
	ByPriority      map[string]int `json:"by_priority"`
	TotalCost       float64        `json:"total_cost"`
}

type MaintenanceReportResponse struct {
	Summary MaintenanceReportSummary `json:"summary"`
	Data    []MaintenanceResponse    `json:"data"`
}

type ReportFilterQuery struct {
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
	SessionID string
	UserID    string
	AssetID   string
	Status    string
	Category  string
	Priority  string
	Format    string // "json" (default) atau "csv"
}
