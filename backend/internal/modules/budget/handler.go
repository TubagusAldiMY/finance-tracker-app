package budget

import "github.com/gofiber/fiber/v2"

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateBudget godoc
// @Summary      Create Monthly Budget
// @Description  Create a new global budget for a specific month
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateBudgetRequest true "Budget Payload"
// @Success      201  {object}  MonthlyBudget
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /budgets [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CreateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	budget, err := h.useCase.CreateBudget(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": budget})
}

// ListBudget godoc
// @Summary      List Budgets
// @Description  Get list of monthly budgets with date filter
// @Tags         Budget
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date query string false "Start Date (YYYY-MM-DD)"
// @Param        end_date   query string false "End Date (YYYY-MM-DD)"
// @Success      200  {object}  []MonthlyBudget
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /budgets [get]
func (h *Handler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	req := &ListBudgetRequest{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	budgets, err := h.useCase.GetBudgets(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if budgets == nil {
		budgets = []*MonthlyBudget{}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": budgets})
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	// Grouping route: /api/budgets
	// Middleware dipasang di level group ini
	api := app.Group("/api/budgets", authMiddleware)

	api.Post("/", h.Create)
	api.Get("/", h.List)
}
