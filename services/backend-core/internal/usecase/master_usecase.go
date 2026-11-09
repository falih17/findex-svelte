package usecase

import (
	"context"
	"errors"
	"fmt"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/aiservice"
	"github.com/google/uuid"
)

type MasterUsecase interface {
	// Company (Owner only)
	ListCompanies(ctx context.Context) ([]domain.Company, error)
	GetCompany(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	CreateCompany(ctx context.Context, c *domain.Company) error
	UpdateCompany(ctx context.Context, c *domain.Company) error

	// Locations / Branches (Owner & Admin Company)
	ListLocations(ctx context.Context, companyID uuid.UUID) ([]domain.Branch, error)
	CreateLocation(ctx context.Context, b *domain.Branch) error
	UpdateLocation(ctx context.Context, b *domain.Branch) error
	DeleteLocation(ctx context.Context, id uuid.UUID) error

	// Employees (Owner & Admin Company)
	ListEmployees(ctx context.Context, companyID uuid.UUID) ([]domain.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	CreateEmployee(ctx context.Context, emp *domain.Employee) error
	UpdateEmployee(ctx context.Context, emp *domain.Employee) error
	DeleteEmployee(ctx context.Context, id uuid.UUID) error

	// Face Biometrics (Bisa diunggah oleh Admin Company untuk pegawai, atau self-enroll)
	EnrollFaceBiometric(ctx context.Context, companyID, employeeID uuid.UUID, photoBytes []byte, filename string) (*domain.EmployeeBiometrics, error)
}

type masterUsecase struct {
	companyRepo domain.CompanyRepository
	branchRepo  domain.BranchRepository
	empRepo     domain.EmployeeRepository
	aiClient    *aiservice.Client
}

func NewMasterUsecase(
	comp domain.CompanyRepository,
	branch domain.BranchRepository,
	emp domain.EmployeeRepository,
	ai *aiservice.Client,
) MasterUsecase {
	return &masterUsecase{
		companyRepo: comp,
		branchRepo:  branch,
		empRepo:     emp,
		aiClient:    ai,
	}
}

// Company
func (u *masterUsecase) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	return u.companyRepo.List(ctx)
}

func (u *masterUsecase) GetCompany(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return u.companyRepo.GetByID(ctx, id)
}

func (u *masterUsecase) CreateCompany(ctx context.Context, c *domain.Company) error {
	return u.companyRepo.Create(ctx, c)
}

func (u *masterUsecase) UpdateCompany(ctx context.Context, c *domain.Company) error {
	return u.companyRepo.Update(ctx, c)
}

// Locations
func (u *masterUsecase) ListLocations(ctx context.Context, companyID uuid.UUID) ([]domain.Branch, error) {
	return u.branchRepo.ListByCompany(ctx, companyID)
}

func (u *masterUsecase) CreateLocation(ctx context.Context, b *domain.Branch) error {
	return u.branchRepo.Create(ctx, b)
}

func (u *masterUsecase) UpdateLocation(ctx context.Context, b *domain.Branch) error {
	return u.branchRepo.Update(ctx, b)
}

func (u *masterUsecase) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	return u.branchRepo.Delete(ctx, id)
}

// Employees
func (u *masterUsecase) ListEmployees(ctx context.Context, companyID uuid.UUID) ([]domain.Employee, error) {
	return u.empRepo.ListByCompany(ctx, companyID)
}

func (u *masterUsecase) GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	return u.empRepo.GetByID(ctx, id)
}

func (u *masterUsecase) CreateEmployee(ctx context.Context, emp *domain.Employee) error {
	return u.empRepo.Create(ctx, emp)
}

func (u *masterUsecase) UpdateEmployee(ctx context.Context, emp *domain.Employee) error {
	return u.empRepo.Update(ctx, emp)
}

func (u *masterUsecase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return u.empRepo.Delete(ctx, id)
}

// Face Biometric Enrollment
func (u *masterUsecase) EnrollFaceBiometric(ctx context.Context, companyID, employeeID uuid.UUID, photoBytes []byte, filename string) (*domain.EmployeeBiometrics, error) {
	emp, err := u.empRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, errors.New("pegawai tidak ditemukan")
	}

	// 1. Send to AI service for feature extraction & anti-spoofing check
	aiRes, err := u.aiClient.ExtractFaceAndVerifyLiveness(ctx, photoBytes, filename)
	if err != nil {
		return nil, fmt.Errorf("biometric extraction error: %w", err)
	}

	if !aiRes.FaceDetected {
		return nil, errors.New("wajah tidak terdeteksi pada foto. Pastikan foto tegak dan jelas.")
	}

	if !aiRes.IsLive {
		return nil, errors.New("foto terindikasi manipulasi/spoofing (kurang tajam atau dari layar ponsel).")
	}

	// 2. Save master biometric
	bio := &domain.EmployeeBiometrics{
		CompanyID:           companyID,
		EmployeeID:          emp.ID,
		MasterImageURL:      fmt.Sprintf("master-faces/%s_%s", emp.ID, filename),
		FaceEmbedding:       aiRes.Embedding,
		LivenessScoreSample: aiRes.LivenessScore,
	}

	if err := u.empRepo.SaveBiometrics(ctx, bio); err != nil {
		return nil, fmt.Errorf("gagal menyimpan basis data biometrik: %w", err)
	}

	return bio, nil
}
