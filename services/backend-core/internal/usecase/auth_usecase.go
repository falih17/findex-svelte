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
	empRepo   domain.EmployeeRepository
	jwtSecret string
}

func NewAuthUsecase(empRepo domain.EmployeeRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		empRepo:   empRepo,
		jwtSecret: jwtSecret,
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
