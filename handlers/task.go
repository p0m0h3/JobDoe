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
// @Success      201 {object} schemas.Task
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

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTask godoc
// @Summary      Get the details of a task
// @Description  Returns the details of a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "task id"
// @Success      200 {object} schemas.Task
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id} [get]
func GetTask(c *fiber.Ctx) error {
	result, ok := state.Tasks[c.Params("id")]
	state.UpdateTask(result)
	if !ok {
		return NotFoundError(c)
	}
	return c.JSON(result)
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task's data and container
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Param        id path string true "task id"
// @Success      204
// @Failure      404 {object} schemas.ErrorResponse
// @Router       /task/{id} [delete]
func DeleteTask(c *fiber.Ctx) error {
	result, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}

	if err := podman.DeleteContainer(result.ID); err != nil {
		return InternalServerError(c)
	}

	return c.Status(fiber.StatusNoContent).SendString("")
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

// GetTaskStats godoc
// @Summary      Task statistics
// @Description  Get the resource usage of a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "task id"
// @Success      200 {object} string
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id}/stats [get]
func GetTaskStats(c *fiber.Ctx) error {
	t, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}
	data, err := podman.GetContainerStats(t.ID)
	if err != nil || len(data.Stats) < 1 {
		return InternalServerError(c)
	}

	stats := data.Stats[0]

	return c.JSON(schemas.GetTaskStatsResponse{
		ID:          stats.ContainerID,
		AvgCPU:      stats.AvgCPU,
		CPU:         stats.CPU,
		MemUsage:    stats.MemUsage,
		MemLimit:    stats.MemLimit,
		MemPerc:     stats.MemPerc,
		NetInput:    stats.NetInput,
		NetOutput:   stats.NetOutput,
		BlockInput:  stats.BlockInput,
		BlockOutput: stats.BlockOutput,
		UpTime:      stats.UpTime,
		Duration:    stats.Duration,
	})
}

// WaitOnTask godoc
// @Summary      Wait on task
// @Description  Return when a task state is changed to exited
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Param        id path string true "task id"
// @Success      200 {object} string
// @Failure      500 {object} schemas.ErrorResponse
// @Failure      404 {object} schemas.ErrorResponse
// @Router       /task/{id}/wait [get]
func WaitOnTask(c *fiber.Ctx) error {
	t, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}

	err := podman.WaitOnContainer(t.ID)
	if err != nil {
		return InternalServerError(c)
	}

	state.UpdateTask(t)

	return c.JSON(t)
}
