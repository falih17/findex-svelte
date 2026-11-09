package http

import (
	"io"
	"strconv"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AttendanceHandler struct {
	attUsecase usecase.AttendanceUsecase
}

func NewAttendanceHandler(u usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{attUsecase: u}
}

func (h *AttendanceHandler) CheckIn(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Foto wajah (photo) wajib diunggah")
	}

	latStr := c.FormValue("latitude")
	lngStr := c.FormValue("longitude")
	empIDStr := c.FormValue("employee_id", "55555555-5555-5555-5555-555555555555")
	branchIDStr := c.FormValue("branch_id", "22222222-2222-2222-2222-222222222222")

	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "employee_id tidak valid")
	}

	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "branch_id tidak valid")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "latitude tidak valid")
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "longitude tidak valid")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal membaca file foto")
	}
	defer file.Close()

	photoBytes, err := io.ReadAll(file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Gagal memproses byte foto")
	}

	req := domain.CheckInRequest{
		EmployeeID: empID,
		BranchID:   branchID,
		Latitude:   lat,
		Longitude:  lng,
		PhotoBytes: photoBytes,
		Filename:   fileHeader.Filename,
	}

	result, err := h.attUsecase.CheckIn(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, result.Message, result)
}

func (h *AttendanceHandler) History(c *fiber.Ctx) error {
	empIDStr := c.Query("employee_id", "55555555-5555-5555-5555-555555555555")
	empID, err := uuid.Parse(empIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "employee_id tidak valid")
	}

	list, err := h.attUsecase.GetRecentHistory(c.Context(), empID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Riwayat presensi berhasil dimuat", list)
}
