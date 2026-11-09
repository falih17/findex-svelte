package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string

const (
	StatusOnTime        AttendanceStatus = "ON_TIME"
	StatusLate          AttendanceStatus = "LATE"
	StatusEarlyLeave    AttendanceStatus = "EARLY_LEAVE"
	StatusWFH           AttendanceStatus = "WFH"
	StatusOutOfGeofence AttendanceStatus = "OUT_OF_GEOFENCE"
)

type Attendance struct {
	ID                     uuid.UUID        `json:"id"`
	CompanyID              uuid.UUID        `json:"company_id"`
	EmployeeID             uuid.UUID        `json:"employee_id"`
	BranchID               *uuid.UUID       `json:"branch_id,omitempty"`
	AttendanceDate         string           `json:"attendance_date"` // YYYY-MM-DD
	CheckInTime            *time.Time       `json:"check_in_time"`
	CheckInPhotoURL        string           `json:"check_in_photo_url,omitempty"`
	CheckInLat             float64          `json:"check_in_lat"`
	CheckInLng             float64          `json:"check_in_lng"`
	CheckInDistanceMeters  float64          `json:"check_in_distance_meters"`
	CheckInSimilarityScore float64          `json:"check_in_similarity_score"`
	CheckInStatus          AttendanceStatus `json:"check_in_status"`
	CheckOutTime           *time.Time       `json:"check_out_time,omitempty"`
	CheckOutStatus         string           `json:"check_out_status,omitempty"`
	WorkDurationMinutes    int              `json:"work_duration_minutes"`
	IsAnomaly              bool             `json:"is_anomaly"`
	Notes                  string           `json:"notes,omitempty"`
}

type CheckInRequest struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	BranchID   uuid.UUID `json:"branch_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	PhotoBytes []byte    `json:"-"`
	Filename   string    `json:"-"`
}

type CheckInResult struct {
	AttendanceID    uuid.UUID        `json:"attendance_id"`
	Status          AttendanceStatus `json:"status"`
	CheckInTime     string           `json:"check_in_time"`
	SimilarityScore float64          `json:"similarity_score"`
	LivenessScore   float64          `json:"liveness_score"`
	DistanceMeters  float64          `json:"distance_meters"`
	InsideGeofence  bool             `json:"inside_geofence"`
	Message         string           `json:"message"`
}

type AttendanceRepository interface {
	RecordCheckIn(ctx context.Context, att *Attendance) error
	RecordCheckOut(ctx context.Context, id uuid.UUID, checkOutTime time.Time, status string, durationMinutes int) error
	GetTodayAttendance(ctx context.Context, employeeID uuid.UUID, date string) (*Attendance, error)
	ListRecent(ctx context.Context, employeeID uuid.UUID, limit int) ([]Attendance, error)
}
