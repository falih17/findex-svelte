package http

import (
	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BranchHandler struct {
	branchUsecase usecase.BranchUsecase
}

func NewBranchHandler(u usecase.BranchUsecase) *BranchHandler {
	return &BranchHandler{branchUsecase: u}
}

func (h *BranchHandler) List(c *fiber.Ctx) error {
	compIDStr := c.Query("company_id", "11111111-1111-1111-1111-111111111111")
	companyID, err := uuid.Parse(compIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "company_id tidak valid")
	}

	branches, err := h.branchUsecase.ListBranches(c.Context(), companyID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Daftar cabang berhasil dimuat", branches)
}

func (h *BranchHandler) Create(c *fiber.Ctx) error {
	var b domain.Branch
	if err := c.BodyParser(&b); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload cabang tidak valid")
	}

	if err := h.branchUsecase.CreateBranch(c.Context(), &b); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Cabang baru berhasil dibuat", b)
}
