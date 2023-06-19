package tool

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAllTools godoc
// @Summary      Get all available tools
// @Description  Get the name of all available tools
// @Tags         tools
// @Accept       json
// @Produce      json
// @Success      200 {array} string
// @Failure      500 {object} ErrorResponse
// @Router       /tool/ [get]
func GetAllTools(c *fiber.Ctx) error {
	result := make([]string, 0)
	for name := range Tools {
		result = append(result, name)
	}
	return c.JSON(result)
}

// GetTool godoc
// @Summary      Get available tool by name
// @Description  Get details for a tool
// @Tags         tools
// @Param        name  path  string  true  "Tool name"
// @Accept       json
// @Produce      json
// @Success      200 {object} GetToolResponse
// @Failure      500 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /tool/{name} [get]
func GetTool(c *fiber.Ctx) error {
	tool, ok := Tools[c.Params("name")]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Code: fiber.StatusNotFound,
		})
	}

	return c.JSON(
		GetToolResponse{
			Name: strings.TrimSuffix(c.Params("name"), ".toml"),
			Spec: *tool,
		},
	)
}
