package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EmployeeRole string

const (
	RoleSuperAdmin EmployeeRole = "SUPERADMIN"
	RoleHRAdmin    EmployeeRole = "HR_ADMIN"
	RoleManager    EmployeeRole = "MANAGER"
	RoleEmployee   EmployeeRole = "EMPLOYEE"
)

type Employee struct {
	ID              uuid.UUID    `json:"id"`
	CompanyID       uuid.UUID    `json:"company_id"`
	DepartmentID    *uuid.UUID   `json:"department_id,omitempty"`
	DefaultBranchID *uuid.UUID   `json:"default_branch_id,omitempty"`
	NIK             string       `json:"nik"`
	Name            string       `json:"name"`
	Email           string       `json:"email"`
	PasswordHash    string       `json:"-"`
	Role            EmployeeRole `json:"role"`
	EmploymentType  string       `json:"employment_type"`
	AllowWFH        bool         `json:"allow_wfh"`
	IsActive        bool         `json:"is_active"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type EmployeeBiometrics struct {
	ID                  uuid.UUID `json:"id"`
	CompanyID           uuid.UUID `json:"company_id"`
	EmployeeID          uuid.UUID `json:"employee_id"`
	MasterImageURL      string    `json:"master_image_url"`
	FaceEmbedding       []float64 `json:"-"` // 512 dimensions
	LivenessScoreSample float64   `json:"liveness_score_sample"`
	RegisteredAt        time.Time `json:"registered_at"`
}

type EmployeeRepository interface {
	GetByEmail(ctx context.Context, email string) (*Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	GetBiometrics(ctx context.Context, employeeID uuid.UUID) (*EmployeeBiometrics, error)
	SaveBiometrics(ctx context.Context, bio *EmployeeBiometrics) error
	CheckFaceSimilarity(ctx context.Context, employeeID uuid.UUID, liveVector []float64) (similarity float64, err error)
}
