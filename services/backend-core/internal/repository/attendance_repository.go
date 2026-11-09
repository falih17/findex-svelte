package repository

import (
	"context"
	"fmt"
	"time"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type attendanceRepository struct {
	db *pgxpool.Pool
}

func NewAttendanceRepository(db *pgxpool.Pool) domain.AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) RecordCheckIn(ctx context.Context, att *domain.Attendance) error {
	query := `
		INSERT INTO attendances (
			id, company_id, employee_id, branch_id, attendance_date,
			check_in_time, check_in_photo_url, check_in_lat, check_in_lng,
			check_in_location_geom, check_in_distance_meters,
			check_in_similarity_score, check_in_status
		) VALUES (
			$1, $2, $3, $4, $5::date,
			$6, $7, $8, $9,
			ST_SetSRID(ST_MakePoint($9, $8), 4326), $10,
			$11, $12
		) RETURNING id`

	if att.ID == uuid.Nil {
		att.ID = uuid.New()
	}
	now := time.Now()
	att.CheckInTime = &now

	return r.db.QueryRow(ctx, query,
		att.ID, att.CompanyID, att.EmployeeID, att.BranchID, att.AttendanceDate,
		att.CheckInTime, att.CheckInPhotoURL, att.CheckInLat, att.CheckInLng,
		att.CheckInDistanceMeters, att.CheckInSimilarityScore, att.CheckInStatus,
	).Scan(&att.ID)
}

func (r *attendanceRepository) GetTodayAttendance(ctx context.Context, employeeID uuid.UUID, date string) (*domain.Attendance, error) {
	query := `
		SELECT 
			id, company_id, employee_id, branch_id, attendance_date,
			check_in_time, check_in_photo_url, check_in_lat, check_in_lng,
			check_in_distance_meters, check_in_similarity_score, check_in_status,
			check_out_time, check_out_status, work_duration_minutes, is_anomaly
		FROM attendances
		WHERE employee_id = $1 AND attendance_date = $2::date
		LIMIT 1`

	var a domain.Attendance
	err := r.db.QueryRow(ctx, query, employeeID, date).Scan(
		&a.ID, &a.CompanyID, &a.EmployeeID, &a.BranchID, &a.AttendanceDate,
		&a.CheckInTime, &a.CheckInPhotoURL, &a.CheckInLat, &a.CheckInLng,
		&a.CheckInDistanceMeters, &a.CheckInSimilarityScore, &a.CheckInStatus,
		&a.CheckOutTime, &a.CheckOutStatus, &a.WorkDurationMinutes, &a.IsAnomaly,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get today attendance: %w", err)
	}
	return &a, nil
}

func (r *attendanceRepository) ListRecent(ctx context.Context, employeeID uuid.UUID, limit int) ([]domain.Attendance, error) {
	query := `
		SELECT 
			id, company_id, employee_id, branch_id, attendance_date,
			check_in_time, check_in_distance_meters, check_in_similarity_score, check_in_status,
			check_out_time, check_out_status, work_duration_minutes
		FROM attendances
		WHERE employee_id = $1
		ORDER BY attendance_date DESC, check_in_time DESC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, employeeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Attendance
	for rows.Next() {
		var a domain.Attendance
		if err := rows.Scan(
			&a.ID, &a.CompanyID, &a.EmployeeID, &a.BranchID, &a.AttendanceDate,
			&a.CheckInTime, &a.CheckInDistanceMeters, &a.CheckInSimilarityScore, &a.CheckInStatus,
			&a.CheckOutTime, &a.CheckOutStatus, &a.WorkDurationMinutes,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	return list, nil
}
