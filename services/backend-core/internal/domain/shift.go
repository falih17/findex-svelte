package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID                   uuid.UUID `json:"id"`
	CompanyID            uuid.UUID `json:"company_id"`
	Name                 string    `json:"name"`                   // Cth: Shift Reguler Pagi, Shift Malam
	StartTime            string    `json:"start_time"`             // 08:00:00 (Jam Masuk Kerja)
	CheckInDeadline      string    `json:"check_in_deadline"`      // 08:30:00 (Batas Maksimal Check In)
	EarliestCheckoutTime string    `json:"earliest_checkout_time"` // 16:30:00 (Waktu Paling Awal Boleh Check Out)
	EndTime              string    `json:"end_time"`               // 17:00:00 (Jam Selesai Kerja)
	GracePeriodLateMin   int       `json:"grace_period_late_minutes"`
	MinWorkHours         float64   `json:"min_work_hours"`
	CreatedAt            time.Time `json:"created_at"`
}

type EmployeeScheduleAssignment struct {
	CompanyID    uuid.UUID `json:"company_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	ShiftID      uuid.UUID `json:"shift_id"`
	ScheduleDate string    `json:"schedule_date"` // YYYY-MM-DD
	IsOffDay     bool      `json:"is_off_day"`
}

type LeaveRequest struct {
	ID            uuid.UUID  `json:"id"`
	CompanyID     uuid.UUID  `json:"company_id"`
	EmployeeID    uuid.UUID  `json:"employee_id"`
	LeaveTypeID   uuid.UUID  `json:"leave_type_id"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	TotalDays     int        `json:"total_days"`
	Reason        string     `json:"reason"`
	AttachmentURL string     `json:"attachment_url,omitempty"`
	Status        string     `json:"status"` // PENDING, APPROVED, REJECTED
	ApprovedBy    *uuid.UUID `json:"approved_by,omitempty"`
	ApprovalNotes string     `json:"approval_notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ShiftRepository interface {
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]Shift, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Shift, error)
	Create(ctx context.Context, s *Shift) error
	Update(ctx context.Context, s *Shift) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignToEmployee(ctx context.Context, assignment *EmployeeScheduleAssignment) error
}

type LeaveRepository interface {
	CreateRequest(ctx context.Context, req *LeaveRequest) error
	ListByEmployee(ctx context.Context, employeeID uuid.UUID) ([]LeaveRequest, error)
}
