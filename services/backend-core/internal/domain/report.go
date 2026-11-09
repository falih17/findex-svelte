package domain

import (
	"context"

	"github.com/google/uuid"
)

type AttendanceReportFilter struct {
	CompanyID  uuid.UUID  `json:"company_id"`
	StartDate  string     `json:"start_date"` // YYYY-MM-DD
	EndDate    string     `json:"end_date"`   // YYYY-MM-DD
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	EmployeeID *uuid.UUID `json:"employee_id,omitempty"`
	Status     string     `json:"status,omitempty"`
}

type AttendanceReportSummary struct {
	TotalRecords    int     `json:"total_records"`
	OnTimeCount     int     `json:"on_time_count"`
	LateCount       int     `json:"late_count"`
	EarlyLeaveCount int     `json:"early_leave_count"`
	OutOfGeofence   int     `json:"out_of_geofence_count"`
	WFHCount        int     `json:"wfh_count"`
	AvgSimilarity   float64 `json:"avg_face_similarity"`
}

type AttendanceReportItem struct {
	ID              uuid.UUID `json:"id"`
	AttendanceDate  string    `json:"attendance_date"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	EmployeeNIK     string    `json:"employee_nik"`
	EmployeeName    string    `json:"employee_name"`
	BranchName      string    `json:"branch_name"`
	CheckInTime     *string   `json:"check_in_time"`
	CheckInStatus   string    `json:"check_in_status"`
	CheckInDistance float64   `json:"check_in_distance_meters"`
	FaceSimilarity  float64   `json:"face_similarity_score"`
	CheckOutTime    *string   `json:"check_out_time,omitempty"`
	CheckOutStatus  string    `json:"check_out_status,omitempty"`
	WorkDurationMin int       `json:"work_duration_minutes"`
}

type AttendanceReportResult struct {
	Summary AttendanceReportSummary `json:"summary"`
	Items   []AttendanceReportItem  `json:"items"`
}

type ReportRepository interface {
	GetAttendanceReport(ctx context.Context, filter AttendanceReportFilter) (*AttendanceReportResult, error)
}
