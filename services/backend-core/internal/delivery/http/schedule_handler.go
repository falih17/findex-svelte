package http

import (
	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ScheduleHandler struct {
	scheduleUsecase usecase.ScheduleUsecase
}

func NewScheduleHandler(u usecase.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{scheduleUsecase: u}
}

func (h *ScheduleHandler) ListShifts(c *fiber.Ctx) error {
	companyID, err := getCompanyIDFromContextOrQuery(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	shifts, err := h.scheduleUsecase.ListShifts(c.Context(), companyID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Daftar jadwal shift berhasil dimuat", shifts)
}

func (h *ScheduleHandler) CreateShift(c *fiber.Ctx) error {
	var s domain.Shift
	if err := c.BodyParser(&s); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload shift tidak valid")
	}

	if s.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		s.CompanyID = compID
	}

	if err := h.scheduleUsecase.CreateShift(c.Context(), &s); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Jadwal shift baru berhasil dibuat", s)
}

func (h *ScheduleHandler) UpdateShift(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID shift tidak valid")
	}

	var s domain.Shift
	if err := c.BodyParser(&s); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload shift tidak valid")
	}
	s.ID = id

	if s.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		s.CompanyID = compID
	}

	if err := h.scheduleUsecase.UpdateShift(c.Context(), &s); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Jadwal shift berhasil diperbarui", s)
}

func (h *ScheduleHandler) DeleteShift(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID shift tidak valid")
	}

	if err := h.scheduleUsecase.DeleteShift(c.Context(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Jadwal shift berhasil dihapus", nil)
}

func (h *ScheduleHandler) AssignShift(c *fiber.Ctx) error {
	var a domain.EmployeeScheduleAssignment
	if err := c.BodyParser(&a); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload assignment jadwal tidak valid")
	}

	if a.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		a.CompanyID = compID
	}

	if err := h.scheduleUsecase.AssignShift(c.Context(), &a); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Jadwal pegawai berhasil ditetapkan", a)
}
