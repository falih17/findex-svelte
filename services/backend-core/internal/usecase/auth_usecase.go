package usecase

import (
	"context"
	"errors"
	"time"

	"findex/backend-core/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*LoginResult, error)
}

type LoginResult struct {
	Token    string          `json:"token"`
	Employee domain.Employee `json:"employee"`
}

type authUsecase struct {
	empRepo     domain.EmployeeRepository
	compRepo    domain.CompanyRepository
	jwtSecret   string
}

func NewAuthUsecase(empRepo domain.EmployeeRepository, compRepo domain.CompanyRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		empRepo:     empRepo,
		compRepo:    compRepo,
		jwtSecret:   jwtSecret,
	}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	emp, err := u.empRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(emp.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Owner can always log in. For ADMIN_COMPANY and EMPLOYEE, verify if company is active
	if emp.Role != domain.RoleOwner && u.compRepo != nil {
		comp, err := u.compRepo.GetByID(ctx, emp.CompanyID)
		if err == nil && comp != nil && !comp.IsActive {
			return nil, errors.New("perusahaan Anda sedang dinonaktifkan oleh administrator. Silakan hubungi pemilik platform.")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"employee_id": emp.ID.String(),
		"company_id":  emp.CompanyID.String(),
		"name":        emp.Name,
		"email":       emp.Email,
		"role":        emp.Role,
		"allow_wfh":   emp.AllowWFH,
		"exp":         time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, errors.New("gagal menerbitkan access token")
	}

	return &LoginResult{
		Token:    tokenStr,
		Employee: *emp,
	}, nil
}
