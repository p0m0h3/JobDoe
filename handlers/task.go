package handlers

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
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

// GetAllTasks godoc
// @Summary      Get a list of all tasks
// @Description  Returns all identified tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200 {object} schemas.Task
// @Router       /task [get]
func GetAllTasks(c *fiber.Ctx) error {
	return c.JSON(state.Tasks)
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

// GetTaskOutputFiles godoc
// @Summary      Get output files from task
// @Description  Returns the contents of output files
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "task id"
// @Success      200 {object} map[string]string
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id}/files [get]
func GetTaskOutputFiles(c *fiber.Ctx) error {
	task, ok := state.Tasks[c.Params("id")]
	state.UpdateTask(task)
	if !ok || task.Status != "exited" {
		return NotFoundError(c)
	}
	archive := &bytes.Buffer{}

	podman.CopyFromContainer(task.ID, archive, fmt.Sprint(state.FILES_PREFIX, state.OUTPUT_PREFIX))

	data := tar.NewReader(archive)

	result := make(map[string]string)

	for {
		hdr, err := data.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return InternalServerError(c)
		}

		output := &strings.Builder{}
		b64 := base64.NewEncoder(base64.StdEncoding, output)
		if _, err := io.Copy(b64, data); err != nil {
			return InternalServerError(c)
		}
		if output.Len() > 0 {
			result[strings.TrimPrefix(hdr.Name, state.OUTPUT_PREFIX)] = output.String()
		}
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

	delete(state.Tasks, result.ID)

	return c.Status(fiber.StatusNoContent).SendString("")
}

// PruneTasks godoc
// @Summary      Prune all stopped/exited tasks
// @Description  Prune the data and container of all stopped/exited tasks
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Success      204
// @Router       /task [delete]
func PruneTasks(c *fiber.Ctx) error {
	podman.PruneTasks()
	state.ResetTasks()
	return c.Status(fiber.StatusNoContent).SendString("")
}

// GetTaskLog godoc
// @Summary      Get the stdout/stderr of a exited task
// @Description  Get the stdout/stderr of an exited task in plaintext
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Param        id path string true "task id"
// @Param        stderr query bool false "should include stderr"
// @Success      200 {object} string
// @Failure      404 {object} schemas.ErrorResponse
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /task/{id}/log [get]
func GetTaskLog(c *fiber.Ctx) error {
	t, ok := state.Tasks[c.Params("id")]
	if !ok {
		return NotFoundError(c)
	}

	stderr, err := strconv.ParseBool(c.Query("stderr", "false"))
	if err != nil {
		return BadRequestError(c, []string{"stderr query parameter couldn't be parsed"})
	}

	output := make(chan string, 1024)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = podman.GetContainerLog(t.ID, stderr, output)
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

func StreamTaskLog(c *websocket.Conn) {
	t, ok := state.Tasks[c.Params("id")]
	if !ok {
		c.WriteMessage(1, []byte(fiber.ErrNotFound.Message))
		c.Close()
	}

	var err error
	output := make(chan string, 1024)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err = podman.StreamContainerLog(t.ID, true, output)
		cancel()
	}()

	if err != nil {
		c.WriteMessage(1, []byte(fiber.ErrNotFound.Message))
		c.Close()
	}

	for {
		select {
		case <-ctx.Done():
			c.Close()
		case line := <-output:
			c.WriteMessage(1, []byte(line))
		}
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
