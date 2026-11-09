package repository

import (
	"context"
	"fmt"
	"strings"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type employeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	query := `
		SELECT id, company_id, department_id, default_branch_id, nik, name, email, 
		       password_hash, role, employment_type, allow_wfh, is_active, created_at, updated_at
		FROM employees
		WHERE email = $1 AND is_active = TRUE
		LIMIT 1`

	var e domain.Employee
	err := r.db.QueryRow(ctx, query, email).Scan(
		&e.ID, &e.CompanyID, &e.DepartmentID, &e.DefaultBranchID,
		&e.NIK, &e.Name, &e.Email, &e.PasswordHash,
		&e.Role, &e.EmploymentType, &e.AllowWFH, &e.IsActive,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, company_id, department_id, default_branch_id, nik, name, email, 
		       password_hash, role, employment_type, allow_wfh, is_active, created_at, updated_at
		FROM employees
		WHERE id = $1
		LIMIT 1`

	var e domain.Employee
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.CompanyID, &e.DepartmentID, &e.DefaultBranchID,
		&e.NIK, &e.Name, &e.Email, &e.PasswordHash,
		&e.Role, &e.EmploymentType, &e.AllowWFH, &e.IsActive,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *employeeRepository) GetBiometrics(ctx context.Context, employeeID uuid.UUID) (*domain.EmployeeBiometrics, error) {
	query := `
		SELECT id, company_id, employee_id, master_image_url, liveness_score_sample, registered_at
		FROM employee_biometrics
		WHERE employee_id = $1
		LIMIT 1`

	var b domain.EmployeeBiometrics
	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&b.ID, &b.CompanyID, &b.EmployeeID,
		&b.MasterImageURL, &b.LivenessScoreSample, &b.RegisteredAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *employeeRepository) SaveBiometrics(ctx context.Context, bio *domain.EmployeeBiometrics) error {
	// Format vector slice into postgres vector format: "[0.123, -0.456, ...]"
	vecStr := formatVector(bio.FaceEmbedding)

	query := `
		INSERT INTO employee_biometrics (company_id, employee_id, master_image_url, face_embedding, liveness_score_sample)
		VALUES ($1, $2, $3, $4::vector, $5)
		ON CONFLICT (employee_id) DO UPDATE SET
			master_image_url = EXCLUDED.master_image_url,
			face_embedding = EXCLUDED.face_embedding,
			liveness_score_sample = EXCLUDED.liveness_score_sample,
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.Exec(ctx, query, bio.CompanyID, bio.EmployeeID, bio.MasterImageURL, vecStr, bio.LivenessScoreSample)
	return err
}

func (r *employeeRepository) CheckFaceSimilarity(ctx context.Context, employeeID uuid.UUID, liveVector []float64) (float64, error) {
	vecStr := formatVector(liveVector)

	// Cosine distance operator <=>: similarity = 1 - (embedding <=> $1)
	query := `
		SELECT 1 - (face_embedding <=> $1::vector) AS similarity
		FROM employee_biometrics
		WHERE employee_id = $2
		LIMIT 1`

	var similarity float64
	err := r.db.QueryRow(ctx, query, vecStr, employeeID).Scan(&similarity)
	if err != nil {
		if err == pgx.ErrNoRows {
			// If employee has not enrolled master face yet, return -1.0
			return -1.0, nil
		}
		return 0.0, fmt.Errorf("vector similarity query failed: %w", err)
	}

	return similarity, nil
}

func formatVector(v []float64) string {
	strs := make([]string, len(v))
	for i, val := range v {
		strs[i] = fmt.Sprintf("%.6f", val)
	}
	return "[" + strings.Join(strs, ",") + "]"
}
