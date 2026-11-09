package repository

import (
	"context"
	"fmt"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type shiftRepository struct {
	db *pgxpool.Pool
}

func NewShiftRepository(db *pgxpool.Pool) domain.ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Shift, error) {
	query := `
		SELECT 
			id, company_id, name, 
			to_char(start_time, 'HH24:MI:SS') AS start_time,
			to_char(check_in_deadline, 'HH24:MI:SS') AS check_in_deadline,
			to_char(earliest_checkout_time, 'HH24:MI:SS') AS earliest_checkout_time,
			to_char(end_time, 'HH24:MI:SS') AS end_time,
			grace_period_late_minutes, min_work_hours, created_at
		FROM shifts
		WHERE company_id = $1
		ORDER BY start_time ASC`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list shifts: %w", err)
	}
	defer rows.Close()

	var list []domain.Shift
	for rows.Next() {
		var s domain.Shift
		if err := rows.Scan(
			&s.ID, &s.CompanyID, &s.Name,
			&s.StartTime, &s.CheckInDeadline, &s.EarliestCheckoutTime, &s.EndTime,
			&s.GracePeriodLateMin, &s.MinWorkHours, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, nil
}

func (r *shiftRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shift, error) {
	query := `
		SELECT 
			id, company_id, name, 
			to_char(start_time, 'HH24:MI:SS') AS start_time,
			to_char(check_in_deadline, 'HH24:MI:SS') AS check_in_deadline,
			to_char(earliest_checkout_time, 'HH24:MI:SS') AS earliest_checkout_time,
			to_char(end_time, 'HH24:MI:SS') AS end_time,
			grace_period_late_minutes, min_work_hours, created_at
		FROM shifts WHERE id = $1`

	var s domain.Shift
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.CompanyID, &s.Name,
		&s.StartTime, &s.CheckInDeadline, &s.EarliestCheckoutTime, &s.EndTime,
		&s.GracePeriodLateMin, &s.MinWorkHours, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *shiftRepository) Create(ctx context.Context, s *domain.Shift) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	query := `
		INSERT INTO shifts (
			id, company_id, name, start_time, check_in_deadline, earliest_checkout_time, end_time,
			grace_period_late_minutes, min_work_hours
		) VALUES ($1, $2, $3, $4::time, $5::time, $6::time, $7::time, $8, $9)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		s.ID, s.CompanyID, s.Name, s.StartTime, s.CheckInDeadline, s.EarliestCheckoutTime, s.EndTime,
		s.GracePeriodLateMin, s.MinWorkHours,
	).Scan(&s.CreatedAt)
}

func (r *shiftRepository) Update(ctx context.Context, s *domain.Shift) error {
	query := `
		UPDATE shifts SET
			name = $1,
			start_time = $2::time,
			check_in_deadline = $3::time,
			earliest_checkout_time = $4::time,
			end_time = $5::time,
			grace_period_late_minutes = $6,
			min_work_hours = $7
		WHERE id = $8 AND company_id = $9`

	_, err := r.db.Exec(ctx, query,
		s.Name, s.StartTime, s.CheckInDeadline, s.EarliestCheckoutTime, s.EndTime,
		s.GracePeriodLateMin, s.MinWorkHours, s.ID, s.CompanyID,
	)
	return err
}

func (r *shiftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shifts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *shiftRepository) AssignToEmployee(ctx context.Context, a *domain.EmployeeScheduleAssignment) error {
	query := `
		INSERT INTO employee_schedules (company_id, employee_id, shift_id, schedule_date, is_off_day)
		VALUES ($1, $2, $3, $4::date, $5)
		ON CONFLICT (company_id, employee_id, schedule_date) DO UPDATE SET
			shift_id = EXCLUDED.shift_id,
			is_off_day = EXCLUDED.is_off_day`

	_, err := r.db.Exec(ctx, query, a.CompanyID, a.EmployeeID, a.ShiftID, a.ScheduleDate, a.IsOffDay)
	return err
}
