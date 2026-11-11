package repository

import (
	"context"
	"fmt"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type companyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) domain.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) List(ctx context.Context) ([]domain.Company, error) {
	query := `
		SELECT 
			c.id, c.code, c.name, c.plan_tier, c.max_employees, c.is_active, 
			c.created_at, c.updated_at,
			s.timezone, s.face_similarity_threshold, s.allow_wfh_default, s.strict_geofence, s.auto_overtime_approval,
			COALESCE(COUNT(DISTINCT e.id), 0) AS total_employees,
			COALESCE(COUNT(DISTINCT b.id), 0) AS total_branches
		FROM companies c
		LEFT JOIN company_settings s ON c.id = s.company_id
		LEFT JOIN employees e ON c.id = e.company_id
		LEFT JOIN branches b ON c.id = b.company_id
		GROUP BY c.id, c.code, c.name, c.plan_tier, c.max_employees, c.is_active, c.created_at, c.updated_at,
		         s.timezone, s.face_similarity_threshold, s.allow_wfh_default, s.strict_geofence, s.auto_overtime_approval
		ORDER BY c.created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()

	var list []domain.Company
	for rows.Next() {
		var c domain.Company
		var tz *string
		var sim *float64
		var wfh, strict, autoOT *bool

		if err := rows.Scan(
			&c.ID, &c.Code, &c.Name, &c.PlanTier, &c.MaxEmployees, &c.IsActive,
			&c.CreatedAt, &c.UpdatedAt,
			&tz, &sim, &wfh, &strict, &autoOT,
			&c.TotalEmployees, &c.TotalBranches,
		); err != nil {
			return nil, err
		}

		if tz != nil {
			c.Settings = &domain.CompanySettings{
				CompanyID:               c.ID,
				Timezone:                *tz,
				FaceSimilarityThreshold: *sim,
				AllowWFHDefault:         *wfh,
				StrictGeofence:          *strict,
				AutoOvertimeApproval:    *autoOT,
			}
		}
		list = append(list, c)
	}

	return list, nil
}

func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	query := `
		SELECT 
			c.id, c.code, c.name, c.plan_tier, c.max_employees, c.is_active, 
			c.created_at, c.updated_at,
			s.timezone, s.face_similarity_threshold, s.allow_wfh_default, s.strict_geofence, s.auto_overtime_approval
		FROM companies c
		LEFT JOIN company_settings s ON c.id = s.company_id
		WHERE c.id = $1`

	var c domain.Company
	var tz *string
	var sim *float64
	var wfh, strict, autoOT *bool

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Code, &c.Name, &c.PlanTier, &c.MaxEmployees, &c.IsActive,
		&c.CreatedAt, &c.UpdatedAt,
		&tz, &sim, &wfh, &strict, &autoOT,
	)
	if err != nil {
		return nil, err
	}

	if tz != nil {
		c.Settings = &domain.CompanySettings{
			CompanyID:               c.ID,
			Timezone:                *tz,
			FaceSimilarityThreshold: *sim,
			AllowWFHDefault:         *wfh,
			StrictGeofence:          *strict,
			AutoOvertimeApproval:    *autoOT,
		}
	}

	return &c, nil
}

func (r *companyRepository) Create(ctx context.Context, c *domain.Company) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	query := `
		INSERT INTO companies (id, code, name, plan_tier, max_employees, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		c.ID, c.Code, c.Name, c.PlanTier, c.MaxEmployees, c.IsActive,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}

	// Create default settings
	settingsQuery := `
		INSERT INTO company_settings (company_id, timezone, face_similarity_threshold, allow_wfh_default, strict_geofence)
		VALUES ($1, 'Asia/Jakarta', 0.80, FALSE, TRUE)
		ON CONFLICT (company_id) DO NOTHING`
	_, _ = r.db.Exec(ctx, settingsQuery, c.ID)

	return nil
}

func (r *companyRepository) Update(ctx context.Context, c *domain.Company) error {
	query := `
		UPDATE companies SET
			name = $1,
			plan_tier = $2,
			max_employees = $3,
			is_active = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query, c.Name, c.PlanTier, c.MaxEmployees, c.IsActive, c.ID)
	return err
}

func (r *companyRepository) SetActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `UPDATE companies SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, isActive, id)
	return err
}

func (r *companyRepository) GetSettings(ctx context.Context, companyID uuid.UUID) (*domain.CompanySettings, error) {
	query := `
		SELECT company_id, timezone, face_similarity_threshold, allow_wfh_default, strict_geofence, auto_overtime_approval
		FROM company_settings WHERE company_id = $1`

	var s domain.CompanySettings
	err := r.db.QueryRow(ctx, query, companyID).Scan(
		&s.CompanyID, &s.Timezone, &s.FaceSimilarityThreshold,
		&s.AllowWFHDefault, &s.StrictGeofence, &s.AutoOvertimeApproval,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *companyRepository) UpdateSettings(ctx context.Context, s *domain.CompanySettings) error {
	query := `
		UPDATE company_settings SET
			timezone = $1,
			face_similarity_threshold = $2,
			allow_wfh_default = $3,
			strict_geofence = $4,
			auto_overtime_approval = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE company_id = $6`

	_, err := r.db.Exec(ctx, query,
		s.Timezone, s.FaceSimilarityThreshold, s.AllowWFHDefault,
		s.StrictGeofence, s.AutoOvertimeApproval, s.CompanyID,
	)
	return err
}
