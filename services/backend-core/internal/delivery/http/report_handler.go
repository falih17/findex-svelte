package http

import (
	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(u usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{reportUsecase: u}
}

func (h *ReportHandler) GetAttendanceReport(c *fiber.Ctx) error {
	companyID, err := getCompanyIDFromContextOrQuery(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	filter := domain.AttendanceReportFilter{
		CompanyID: companyID,
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Status:    c.Query("status"),
	}

	if branchStr := c.Query("branch_id"); branchStr != "" {
		if bid, err := uuid.Parse(branchStr); err == nil {
			filter.BranchID = &bid
		}
	}

	if empStr := c.Query("employee_id"); empStr != "" {
		if eid, err := uuid.Parse(empStr); err == nil {
			filter.EmployeeID = &eid
		}
	}

	report, err := h.reportUsecase.GetAttendanceReport(c.Context(), filter)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Laporan presensi berhasil di-generate", report)
}
