package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	PlanTier     string    `json:"plan_tier"`
	MaxEmployees int       `json:"max_employees"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CompanySettings struct {
	CompanyID                uuid.UUID `json:"company_id"`
	Timezone                 string    `json:"timezone"`
	FaceSimilarityThreshold  float64   `json:"face_similarity_threshold"`
	AllowWFHDefault          bool      `json:"allow_wfh_default"`
	StrictGeofence           bool      `json:"strict_geofence"`
	AutoOvertimeApproval     bool      `json:"auto_overtime_approval"`
}

type CompanyRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetSettings(ctx context.Context, companyID uuid.UUID) (*CompanySettings, error)
}
