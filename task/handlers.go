package task

import (
	"fuzz.codes/fuzzercloud/workerengine/container"
	"github.com/gofiber/contrib/websocket"
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

	badFields, err := ValidateRequest[CreateTaskRequest](req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:       fiber.StatusBadRequest,
			Validation: badFields,
		})
	}

	task, err := NewTask(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	id, err := task.Start()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(CreateTaskResponse{ID: id, Tool: req.ToolName})
}

// GetTask godoc
// @Summary      Get the details of a task
// @Description  Returns the details of a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "task id"
// @Success      200 {object} GetTaskResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /task/{id} [get]
func GetTask(c *fiber.Ctx) error {
	task, err := container.InspectContainer(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return c.JSON(GetTaskResponse{
		ID:        task.ID,
		ImageName: task.ImageName,
		Command:   task.Config.Cmd,
	})
}

func GetTaskOutput(c *websocket.Conn) {
	output, err := container.GetContainerLog(c.Params("id"))
	if err != nil {
		c.WriteMessage(1, []byte(fiber.ErrNotFound.Message))
	}

	for frame := range output {
		c.WriteMessage(1, []byte(frame))
	}
}
