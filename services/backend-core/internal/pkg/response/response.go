package response

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, statusCode int, errMsg string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Error:   errMsg,
	})
}
