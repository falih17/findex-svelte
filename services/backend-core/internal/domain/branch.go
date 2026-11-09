package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID           uuid.UUID `json:"id"`
	CompanyID    uuid.UUID `json:"company_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters int       `json:"radius_meters"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BranchGeofenceCheck struct {
	BranchID       uuid.UUID `json:"branch_id"`
	BranchName     string    `json:"branch_name"`
	IsInside       bool      `json:"is_inside"`
	DistanceMeters float64   `json:"distance_meters"`
	AllowedRadius  int       `json:"allowed_radius"`
}

type BranchRepository interface {
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]Branch, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Branch, error)
	Create(ctx context.Context, b *Branch) error
	CheckGeofence(ctx context.Context, branchID uuid.UUID, lat, lng float64) (*BranchGeofenceCheck, error)
}
