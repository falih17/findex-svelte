package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID                    uuid.UUID `json:"id"`
	CompanyID             uuid.UUID `json:"company_id"`
	Name                  string    `json:"name"`
	StartTime             string    `json:"start_time"` // HH:MM:SS
	EndTime               string    `json:"end_time"`   // HH:MM:SS
	GracePeriodLateMin    int       `json:"grace_period_late_minutes"`
	MinWorkHours          float64   `json:"min_work_hours"`
}

type LeaveRequest struct {
	ID            uuid.UUID `json:"id"`
	CompanyID     uuid.UUID `json:"company_id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LeaveTypeID   uuid.UUID `json:"leave_type_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	TotalDays     int       `json:"total_days"`
	Reason        string    `json:"reason"`
	AttachmentURL string    `json:"attachment_url,omitempty"`
	Status        string    `json:"status"` // PENDING, APPROVED, REJECTED
	ApprovedBy    *uuid.UUID `json:"approved_by,omitempty"`
	ApprovalNotes string    `json:"approval_notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type ShiftRepository interface {
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]Shift, error)
}

type LeaveRepository interface {
	CreateRequest(ctx context.Context, req *LeaveRequest) error
	ListByEmployee(ctx context.Context, employeeID uuid.UUID) ([]LeaveRequest, error)
}
