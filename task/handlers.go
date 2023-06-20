package task

import (
	"fuzz.codes/fuzzercloud/workerengine/tool"
	"github.com/gofiber/fiber/v2"
)

// CreateTask godoc
// @Summary      Create a new task
// @Description  Start a new sandbox with a tool running inside
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body CreateTaskRequest true "new task data"
// @Success      200 {object} CreateTaskResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /task/ [post]
func CreateTask(c *fiber.Ctx) error {
	req := CreateTaskRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusBadRequest})
	}

	badFields, err := ValidateCreateTaskRequest(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:       fiber.StatusBadRequest,
			Validation: badFields,
		})
	}

	spec, err := NewTaskSpec(*tool.Tools[req.ToolName], req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	id, err := NewContainerTask(spec)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code: fiber.StatusInternalServerError,
		})
	}

	return c.JSON(CreateTaskResponse{ID: id, Tool: req.ToolName})
}
