package http

import (
	"findex/backend-core/internal/pkg/response"
	"findex/backend-core/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: u}
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload loginPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Format request tidak valid")
	}

	result, err := h.authUsecase.Login(c.Context(), payload.Email, payload.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Login berhasil", result)
}
