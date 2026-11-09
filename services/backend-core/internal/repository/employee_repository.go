package repository

import (
	"context"
	"fmt"
	"strings"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type employeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Employee, error) {
	query := `
		SELECT 
			e.id, e.company_id, e.department_id, e.default_branch_id, 
			COALESCE(b.name, '') AS branch_name,
			e.nik, e.name, e.email, e.password_hash, e.role, e.employment_type, 
			e.allow_wfh, e.is_active, 
			(eb.id IS NOT NULL) AS has_face_data,
			e.created_at, e.updated_at
		FROM employees e
		LEFT JOIN branches b ON e.default_branch_id = b.id
		LEFT JOIN employee_biometrics eb ON e.id = eb.employee_id
		WHERE e.company_id = $1
		ORDER BY e.name ASC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query employees: %w", err)
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var e domain.Employee
		if err := rows.Scan(
			&e.ID, &e.CompanyID, &e.DepartmentID, &e.DefaultBranchID,
			&e.BranchName,
			&e.NIK, &e.Name, &e.Email, &e.PasswordHash,
			&e.Role, &e.EmploymentType, &e.AllowWFH, &e.IsActive,
			&e.HasFaceData,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	return employees, nil
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
		SELECT e.id, e.company_id, e.department_id, e.default_branch_id, 
		       COALESCE(b.name, ''),
		       e.nik, e.name, e.email, e.password_hash, e.role, e.employment_type, 
		       e.allow_wfh, e.is_active, (eb.id IS NOT NULL), e.created_at, e.updated_at
		FROM employees e
		LEFT JOIN branches b ON e.default_branch_id = b.id
		LEFT JOIN employee_biometrics eb ON e.id = eb.employee_id
		WHERE e.id = $1
		LIMIT 1`

	var e domain.Employee
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.CompanyID, &e.DepartmentID, &e.DefaultBranchID,
		&e.BranchName,
		&e.NIK, &e.Name, &e.Email, &e.PasswordHash,
		&e.Role, &e.EmploymentType, &e.AllowWFH, &e.IsActive,
		&e.HasFaceData,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *employeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	if emp.ID == uuid.Nil {
		emp.ID = uuid.New()
	}

	hash := emp.PasswordHash
	if emp.Password != "" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hash = string(bytes)
	}

	query := `
		INSERT INTO employees (
			id, company_id, department_id, default_branch_id, 
			nik, name, email, password_hash, role, employment_type, allow_wfh, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		emp.ID, emp.CompanyID, emp.DepartmentID, emp.DefaultBranchID,
		emp.NIK, emp.Name, emp.Email, hash, emp.Role, emp.EmploymentType,
		emp.AllowWFH, emp.IsActive,
	).Scan(&emp.CreatedAt, &emp.UpdatedAt)
}

func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	query := `
		UPDATE employees SET
			name = $1,
			nik = $2,
			role = $3,
			employment_type = $4,
			allow_wfh = $5,
			default_branch_id = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8 AND company_id = $9`

	_, err := r.db.Exec(ctx, query,
		emp.Name, emp.NIK, emp.Role, emp.EmploymentType,
		emp.AllowWFH, emp.DefaultBranchID, emp.IsActive,
		emp.ID, emp.CompanyID,
	)
	return err
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM employees WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *employeeRepository) GetBiometrics(ctx context.Context, employeeID uuid.UUID) (*domain.EmployeeBiometrics, error) {
	query := `
		SELECT id, company_id, employee_id, master_image_url, liveness_score_sample, registered_at, updated_at
		FROM employee_biometrics
		WHERE employee_id = $1
		LIMIT 1`

	var b domain.EmployeeBiometrics
	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&b.ID, &b.CompanyID, &b.EmployeeID,
		&b.MasterImageURL, &b.LivenessScoreSample, &b.RegisteredAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *employeeRepository) SaveBiometrics(ctx context.Context, bio *domain.EmployeeBiometrics) error {
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
