package handlers

import (
	"context"
	"strings"

	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/schemas"
	"fuzz.codes/fuzzercloud/workerengine/state"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// CreateTask godoc
// @Summary      Create a new task
// @Description  Start a new sandbox with a tool running inside
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body schemas.CreateTaskRequest true "new task data"
// @Success      200 {object} task.Task
// @Failure      400 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/ [post]
func CreateTask(c *fiber.Ctx) error {
	req := schemas.CreateTaskRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{Code: fiber.StatusBadRequest})
	}

	badFields, err := ValidateRequest[schemas.CreateTaskRequest](req)
	if err != nil {
		return BadRequestError(c, badFields)
	}

	task, err := state.NewTask(req)
	if err != nil {
		return BadRequestError(c, nil)
	}

	id, err := state.StartTask(task)
	if err != nil {
		return InternalServerError(c)
	}

	task.ID = id

	return c.JSON(task)
}

// GetTask godoc
// @Summary      Get the details of a task
// @Description  Returns the details of a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "task id"
// @Success      200 {object} task.Task
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id} [get]
func GetTask(c *fiber.Ctx) error {
	result, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}
	return c.JSON(result)
}

// GetTaskOutput godoc
// @Summary      Get the details of a task
// @Description  Returns the details of a task
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Param        id path string true "task id"
// @Success      200 {object} string
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id}/stdout [get]
func GetTaskOutput(c *fiber.Ctx) error {
	t, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}

	var err error
	output := make(chan string, 1024)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = podman.GetContainerLog(t.ID, output)
		cancel()
	}()

	if err != nil {
		return NotFoundError(c)
	}

	var logs strings.Builder
	for {
		select {
		case <-ctx.Done():
			return c.SendString(logs.String())
		case line := <-output:
			logs.WriteString(line)
		}
	}
}

func StreamTaskOutput(c *websocket.Conn) {
	output := make(chan string, 1024)
	err := podman.GetContainerLog(c.Params("id"), output)
	if err != nil {
		c.WriteMessage(1, []byte(fiber.ErrNotFound.Message))
	}

	for frame := range output {
		c.WriteMessage(1, []byte(frame))
	}
}
