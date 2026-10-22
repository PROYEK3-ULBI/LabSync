package model

type StudentDashboardResponse struct {
	User              UserResponse         `json:"user"`
	FaceEnrolled      bool                 `json:"face_enrolled"`
	ActiveSession     *SessionResponse     `json:"active_session,omitempty"`
	UpcomingSessions  []SessionResponse    `json:"upcoming_sessions"`
	RecentAttendance  []AttendanceResponse `json:"recent_attendance"`
	TotalAttended     int                  `json:"total_attended"`
	CurrentAllocation *AllocationResponse  `json:"current_allocation,omitempty"`
	TotalAllocations  int                  `json:"total_allocations"`
}

type AdminDashboardResponse struct {
	ActiveSessionsCount     int                `json:"active_sessions_count"`
	ActiveSessions          []SessionResponse  `json:"active_sessions"`
	TotalUsers              int                `json:"total_users"`
	TotalStudents           int                `json:"total_students"`
	AssetStatusCounts       map[string]int     `json:"asset_status_counts"`
	ActiveAllocationsCount  int                `json:"active_allocations_count"`
	PendingMaintenanceCount int                `json:"pending_maintenance_count"`
	TotalAttendanceCount    int                `json:"total_attendance_count"`
	RecentAuditLogs         []AuditLogResponse `json:"recent_audit_logs"`
}
