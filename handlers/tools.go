package handlers

import (
	"strings"

	"fuzz.codes/fuzzercloud/workerengine/schemas"
	"fuzz.codes/fuzzercloud/workerengine/state"
	"github.com/gofiber/fiber/v2"
)

// GetAllTools godoc
// @Summary      Get all available tools
// @Description  Get the name of all available tools
// @Security     ApiKeyAuth
// @Tags         tools
// @Accept       json
// @Produce      json
// @Success      200 {array} string
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /v1/tool/ [get]
func GetAllTools(c *fiber.Ctx) error {
	result := make([]string, 0)
	for name := range state.Tools {
		result = append(result, name)
	}
	return c.JSON(result)
}

// GetTool godoc
// @Summary      Get available tool by name
// @Description  Get details for a tool
// @Security     ApiKeyAuth
// @Tags         tools
// @Param        name  path  string  true  "Tool name"
// @Accept       json
// @Produce      json
// @Success      200 {object} schemas.Tool
// @Failure      500 {object} schemas.ErrorResponse
// @Failure      404 {object} schemas.ErrorResponse
// @Router       /v1/tool/{name} [get]
func GetTool(c *fiber.Ctx) error {
	tool, ok := state.Tools[c.Params("name")]
	if !ok {
		return NotFoundError(c)
	}

	return c.JSON(
		schemas.Tool{
			ID:   strings.TrimSuffix(c.Params("name"), ".toml"),
			Spec: tool,
		},
	)
}
