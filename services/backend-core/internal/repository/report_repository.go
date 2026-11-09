package repository

import (
	"context"
	"fmt"
	"strings"

	"findex/backend-core/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) domain.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetAttendanceReport(ctx context.Context, filter domain.AttendanceReportFilter) (*domain.AttendanceReportResult, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, fmt.Sprintf("a.company_id = $%d", argIdx))
	args = append(args, filter.CompanyID)
	argIdx++

	if filter.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("a.attendance_date >= $%d::date", argIdx))
		args = append(args, filter.StartDate)
		argIdx++
	}

	if filter.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("a.attendance_date <= $%d::date", argIdx))
		args = append(args, filter.EndDate)
		argIdx++
	}

	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("a.branch_id = $%d", argIdx))
		args = append(args, *filter.BranchID)
		argIdx++
	}

	if filter.EmployeeID != nil {
		conditions = append(conditions, fmt.Sprintf("a.employee_id = $%d", argIdx))
		args = append(args, *filter.EmployeeID)
		argIdx++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("a.check_in_status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.attendance_date::text,
			a.employee_id,
			e.nik,
			e.name,
			COALESCE(b.name, 'Cabang Fleksibel') AS branch_name,
			to_char(a.check_in_time, 'HH24:MI:SS') AS check_in_time,
			a.check_in_status,
			COALESCE(a.check_in_distance_meters, 0),
			COALESCE(a.check_in_similarity_score, 0),
			to_char(a.check_out_time, 'HH24:MI:SS') AS check_out_time,
			COALESCE(a.check_out_status, ''),
			COALESCE(a.work_duration_minutes, 0)
		FROM attendances a
		JOIN employees e ON a.employee_id = e.id
		LEFT JOIN branches b ON a.branch_id = b.id
		WHERE %s
		ORDER BY a.attendance_date DESC, a.check_in_time DESC`, whereClause)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance report: %w", err)
	}
	defer rows.Close()

	var items []domain.AttendanceReportItem
	var summary domain.AttendanceReportSummary
	var totalSim float64

	for rows.Next() {
		var item domain.AttendanceReportItem
		var checkIn, checkOut *string
		if err := rows.Scan(
			&item.ID, &item.AttendanceDate, &item.EmployeeID,
			&item.EmployeeNIK, &item.EmployeeName, &item.BranchName,
			&checkIn, &item.CheckInStatus, &item.CheckInDistance,
			&item.FaceSimilarity, &checkOut, &item.CheckOutStatus,
			&item.WorkDurationMin,
		); err != nil {
			return nil, err
		}

		item.CheckInTime = checkIn
		item.CheckOutTime = checkOut
		items = append(items, item)

		// Calculate summary
		summary.TotalRecords++
		totalSim += item.FaceSimilarity
		switch item.CheckInStatus {
		case "ON_TIME":
			summary.OnTimeCount++
		case "LATE":
			summary.LateCount++
		case "OUT_OF_GEOFENCE":
			summary.OutOfGeofence++
		case "WFH":
			summary.WFHCount++
		}
		if item.CheckOutStatus == "EARLY_LEAVE" {
			summary.EarlyLeaveCount++
		}
	}

	if summary.TotalRecords > 0 {
		summary.AvgSimilarity = totalSim / float64(summary.TotalRecords)
	}

	return &domain.AttendanceReportResult{
		Summary: summary,
		Items:   items,
	}, nil
}
