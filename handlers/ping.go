package handlers

import (
	"os"

	"git.fuzz.codes/fuzzercloud/workerengine/schemas"
	"git.fuzz.codes/fuzzercloud/workerengine/state"
	"github.com/gofiber/fiber/v2"
)

// Ping godoc
// @Summary      Get info about the worker server
// @Description  Information includes version, configuration and state
// @Security     ApiKeyAuth
// @Tags         state
// @Accept       json
// @Produce      json
// @Success      200 {array} string
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /v1 [get]
func Ping(c *fiber.Ctx) error {
	return c.JSON(schemas.PingResponse{
		Version: "v0.1.0",
		Spec:    "v0.5.3",
		Mode:    os.Getenv("MODE"),
		Tasks:   len(state.Tasks),
		Tools:   len(state.Tools),
	})
}
