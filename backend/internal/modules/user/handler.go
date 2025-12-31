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

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// 1. Parsing Body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request Body",
		})
	}

	// 2. Panggil Usecase
	resp, err := h.useCase.Login(c.Context(), &req)
	if err != nil {
		// 3. Error Handling Spesifik
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Jika Error Lain (Misal DB mati/config error), return 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	// Sukses - Return Token (Status 200 OK)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
	})
}

// Fungsi untuk mendaftarkan route module ini
func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/users")
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
}
