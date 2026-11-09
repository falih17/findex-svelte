package usecase

import (
	"context"

	"findex/backend-core/internal/domain"
	"github.com/google/uuid"
)

type BranchUsecase interface {
	ListBranches(ctx context.Context, companyID uuid.UUID) ([]domain.Branch, error)
	CreateBranch(ctx context.Context, branch *domain.Branch) error
}

type branchUsecase struct {
	repo domain.BranchRepository
}

func NewBranchUsecase(repo domain.BranchRepository) BranchUsecase {
	return &branchUsecase{repo: repo}
}

func (u *branchUsecase) ListBranches(ctx context.Context, companyID uuid.UUID) ([]domain.Branch, error) {
	return u.repo.ListByCompany(ctx, companyID)
}

func (u *branchUsecase) CreateBranch(ctx context.Context, branch *domain.Branch) error {
	return u.repo.Create(ctx, branch)
}
