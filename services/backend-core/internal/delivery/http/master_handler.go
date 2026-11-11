package http

import (
	"fmt"
	"io"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MasterHandler struct {
	masterUsecase usecase.MasterUsecase
}

func NewMasterHandler(u usecase.MasterUsecase) *MasterHandler {
	return &MasterHandler{masterUsecase: u}
}

// -----------------------------------------------------------------------------
// 1. DATA COMPANY (Owner only)
// -----------------------------------------------------------------------------

func (h *MasterHandler) ListCompanies(c *fiber.Ctx) error {
	list, err := h.masterUsecase.ListCompanies(c.Context())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Daftar perusahaan berhasil dimuat", list)
}

func (h *MasterHandler) GetCompany(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID company tidak valid")
	}

	comp, err := h.masterUsecase.GetCompany(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Perusahaan tidak ditemukan")
	}

	return response.Success(c, fiber.StatusOK, "Detail perusahaan", comp)
}

func (h *MasterHandler) CreateCompany(c *fiber.Ctx) error {
	var comp domain.Company
	if err := c.BodyParser(&comp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload company tidak valid")
	}

	if err := h.masterUsecase.CreateCompany(c.Context(), &comp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Perusahaan baru berhasil didaftarkan", comp)
}

func (h *MasterHandler) UpdateCompany(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID company tidak valid")
	}

	var comp domain.Company
	if err := c.BodyParser(&comp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload company tidak valid")
	}
	comp.ID = id

	if err := h.masterUsecase.UpdateCompany(c.Context(), &comp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Data perusahaan berhasil diperbarui", comp)
}

func (h *MasterHandler) ToggleCompanyStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID company tidak valid")
	}

	type StatusRequest struct {
		IsActive bool `json:"is_active"`
	}
	var req StatusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload status tidak valid")
	}

	comp, err := h.masterUsecase.ToggleCompanyActive(c.Context(), id, req.IsActive)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	statusStr := "dinonaktifkan"
	if req.IsActive {
		statusStr = "diaktifkan"
	}

	return response.Success(c, fiber.StatusOK, fmt.Sprintf("Perusahaan '%s' berhasil %s", comp.Name, statusStr), comp)
}

// -----------------------------------------------------------------------------
// 2. DATA LOCATION / GEOFENCE (Owner & Admin Company)
// -----------------------------------------------------------------------------

func (h *MasterHandler) ListLocations(c *fiber.Ctx) error {
	companyID, err := getCompanyIDFromContextOrQuery(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	locations, err := h.masterUsecase.ListLocations(c.Context(), companyID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Daftar lokasi kantor/geofence berhasil dimuat", locations)
}

func (h *MasterHandler) CreateLocation(c *fiber.Ctx) error {
	var b domain.Branch
	if err := c.BodyParser(&b); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload lokasi tidak valid")
	}

	// Default company_id from session if not specified
	if b.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		b.CompanyID = compID
	}

	if err := h.masterUsecase.CreateLocation(c.Context(), &b); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Lokasi absen baru berhasil disimpan", b)
}

func (h *MasterHandler) UpdateLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID lokasi tidak valid")
	}

	var b domain.Branch
	if err := c.BodyParser(&b); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload lokasi tidak valid")
	}
	b.ID = id

	if b.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		b.CompanyID = compID
	}

	if err := h.masterUsecase.UpdateLocation(c.Context(), &b); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Lokasi absen berhasil diperbarui", b)
}

func (h *MasterHandler) DeleteLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID lokasi tidak valid")
	}

	if err := h.masterUsecase.DeleteLocation(c.Context(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Lokasi absen berhasil dihapus", nil)
}

// -----------------------------------------------------------------------------
// 3. DATA EMPLOYEE & FACE BIOMETRIC ENROLLMENT (Owner & Admin Company)
// -----------------------------------------------------------------------------

func (h *MasterHandler) ListEmployees(c *fiber.Ctx) error {
	companyID, err := getCompanyIDFromContextOrQuery(c)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	employees, err := h.masterUsecase.ListEmployees(c.Context(), companyID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Daftar pegawai berhasil dimuat", employees)
}

func (h *MasterHandler) GetEmployee(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	emp, err := h.masterUsecase.GetEmployee(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Pegawai tidak ditemukan")
	}

	return response.Success(c, fiber.StatusOK, "Detail pegawai", emp)
}

func (h *MasterHandler) CreateEmployee(c *fiber.Ctx) error {
	var emp domain.Employee
	if err := c.BodyParser(&emp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload pegawai tidak valid")
	}

	if emp.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		emp.CompanyID = compID
	}

	if emp.Role == "" {
		emp.Role = domain.RoleEmployee
	}

	if err := h.masterUsecase.CreateEmployee(c.Context(), &emp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Pegawai baru berhasil ditambahkan", emp)
}

func (h *MasterHandler) UpdateEmployee(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	var emp domain.Employee
	if err := c.BodyParser(&emp); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Payload pegawai tidak valid")
	}
	emp.ID = id

	if emp.CompanyID == uuid.Nil {
		compID, _ := getCompanyIDFromContextOrQuery(c)
		emp.CompanyID = compID
	}

	if err := h.masterUsecase.UpdateEmployee(c.Context(), &emp); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Data pegawai berhasil diperbarui", emp)
}

func (h *MasterHandler) DeleteEmployee(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	if err := h.masterUsecase.DeleteEmployee(c.Context(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Pegawai berhasil dihapus", nil)
}

// Enroll Face Biometric by Admin Company
func (h *MasterHandler) EnrollFaceBiometric(c *fiber.Ctx) error {
	empID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Foto wajah (photo) wajib diunggah")
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

	compID, _ := getCompanyIDFromContextOrQuery(c)

	bio, err := h.masterUsecase.EnrollFaceBiometric(c.Context(), compID, empID, photoBytes, fileHeader.Filename)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Basis data wajah pegawai berhasil didaftarkan oleh admin", fiber.Map{
		"employee_id":    bio.EmployeeID,
		"liveness_score": bio.LivenessScoreSample,
		"registered_at":  bio.RegisteredAt,
	})
}

func getCompanyIDFromContextOrQuery(c *fiber.Ctx) (uuid.UUID, error) {
	// From JWT context
	if compVal := c.Locals("company_id"); compVal != nil {
		if id, err := uuid.Parse(fmt.Sprintf("%v", compVal)); err == nil && id != uuid.Nil {
			return id, nil
		}
	}
	// Fallback from query parameter
	compQuery := c.Query("company_id", "11111111-1111-1111-1111-111111111111")
	return uuid.Parse(compQuery)
}
