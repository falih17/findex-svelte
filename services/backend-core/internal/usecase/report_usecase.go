package usecase

import (
	"context"

	"findex/backend-core/internal/domain"
)

type ReportUsecase interface {
	GetAttendanceReport(ctx context.Context, filter domain.AttendanceReportFilter) (*domain.AttendanceReportResult, error)
}

type reportUsecase struct {
	repo domain.ReportRepository
}

func NewReportUsecase(repo domain.ReportRepository) ReportUsecase {
	return &reportUsecase{repo: repo}
}

func (u *reportUsecase) GetAttendanceReport(ctx context.Context, filter domain.AttendanceReportFilter) (*domain.AttendanceReportResult, error) {
	return u.repo.GetAttendanceReport(ctx, filter)
}
