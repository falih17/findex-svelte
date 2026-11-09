package usecase

import (
	"context"
	"errors"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
)

type ScheduleUsecase interface {
	ListShifts(ctx context.Context, companyID uuid.UUID) ([]domain.Shift, error)
	CreateShift(ctx context.Context, s *domain.Shift) error
	UpdateShift(ctx context.Context, s *domain.Shift) error
	DeleteShift(ctx context.Context, id uuid.UUID) error
	AssignShift(ctx context.Context, assignment *domain.EmployeeScheduleAssignment) error
}

type scheduleUsecase struct {
	shiftRepo domain.ShiftRepository
}

func NewScheduleUsecase(shiftRepo domain.ShiftRepository) ScheduleUsecase {
	return &scheduleUsecase{shiftRepo: shiftRepo}
}

func (u *scheduleUsecase) ListShifts(ctx context.Context, companyID uuid.UUID) ([]domain.Shift, error) {
	return u.shiftRepo.ListByCompany(ctx, companyID)
}

func (u *scheduleUsecase) CreateShift(ctx context.Context, s *domain.Shift) error {
	if s.Name == "" {
		return errors.New("nama shift wajib diisi")
	}
	if s.StartTime == "" || s.EndTime == "" {
		return errors.New("jam mulai dan jam selesai shift wajib diisi")
	}
	if s.CheckInDeadline == "" {
		s.CheckInDeadline = s.StartTime // Default same as start time
	}
	if s.EarliestCheckoutTime == "" {
		s.EarliestCheckoutTime = s.EndTime // Default same as end time
	}

	return u.shiftRepo.Create(ctx, s)
}

func (u *scheduleUsecase) UpdateShift(ctx context.Context, s *domain.Shift) error {
	return u.shiftRepo.Update(ctx, s)
}

func (u *scheduleUsecase) DeleteShift(ctx context.Context, id uuid.UUID) error {
	return u.shiftRepo.Delete(ctx, id)
}

func (u *scheduleUsecase) AssignShift(ctx context.Context, assignment *domain.EmployeeScheduleAssignment) error {
	return u.shiftRepo.AssignToEmployee(ctx, assignment)
}
