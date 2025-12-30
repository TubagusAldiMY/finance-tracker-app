package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.useCase.Register(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": user})
}

// Fungsi untuk mendaftarkan route module ini
func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/users")
	api.Post("/register", h.Register)
}
