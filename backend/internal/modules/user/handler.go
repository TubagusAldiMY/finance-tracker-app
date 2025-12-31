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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.useCase.Register(c.Context(), &req)
	if err != nil {
		// 1. Cek Error Spesifik (409 Conflict)
		if err == ErrEmailTaken || err == ErrUsernameTaken {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		// 2. Default Error (500 Internal Server Error)
		// Tidak perlu 'if err != nil' lagi disini, karena sudah pasti error (else logic)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// 3. Sukses (201 Created)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})
}

// Fungsi untuk mendaftarkan route module ini
func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/users")
	api.Post("/register", h.Register)
}
