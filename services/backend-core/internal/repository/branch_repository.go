package repository

import (
	"context"
	"fmt"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type branchRepository struct {
	db *pgxpool.Pool
}

func NewBranchRepository(db *pgxpool.Pool) domain.BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Branch, error) {
	query := `
		SELECT id, company_id, name, address, latitude, longitude, radius_meters, created_at, updated_at
		FROM branches
		WHERE company_id = $1
		ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query branches: %w", err)
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var b domain.Branch
		if err := rows.Scan(
			&b.ID, &b.CompanyID, &b.Name, &b.Address,
			&b.Latitude, &b.Longitude, &b.RadiusMeters,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}

	return branches, nil
}

func (r *branchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Branch, error) {
	query := `
		SELECT id, company_id, name, address, latitude, longitude, radius_meters, created_at, updated_at
		FROM branches
		WHERE id = $1`

	var b domain.Branch
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.CompanyID, &b.Name, &b.Address,
		&b.Latitude, &b.Longitude, &b.RadiusMeters,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *branchRepository) Create(ctx context.Context, b *domain.Branch) error {
	query := `
		INSERT INTO branches (id, company_id, name, address, latitude, longitude, radius_meters)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`

	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return r.db.QueryRow(ctx, query,
		b.ID, b.CompanyID, b.Name, b.Address,
		b.Latitude, b.Longitude, b.RadiusMeters,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
}

func (r *branchRepository) CheckGeofence(ctx context.Context, branchID uuid.UUID, lat, lng float64) (*domain.BranchGeofenceCheck, error) {
	query := `
		SELECT 
			id,
			name,
			radius_meters,
			ST_DWithin(location_geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, radius_meters) AS inside,
			ST_Distance(location_geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance
		FROM branches 
		WHERE id = $3 LIMIT 1`

	var check domain.BranchGeofenceCheck
	err := r.db.QueryRow(ctx, query, lng, lat, branchID).Scan(
		&check.BranchID,
		&check.BranchName,
		&check.AllowedRadius,
		&check.IsInside,
		&check.DistanceMeters,
	)
	if err != nil {
		return nil, fmt.Errorf("geofence spatial query failed: %w", err)
	}

	return &check, nil
}
