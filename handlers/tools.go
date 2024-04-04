package handlers

import (
	"encoding/json"

	"git.fuzz.codes/fuzzercloud/tsf"
	"git.fuzz.codes/fuzzercloud/workerengine/state"
	"github.com/gofiber/fiber/v2"
)

// GetAllTools godoc
// @Summary      Get all available tools
// @Description  Get the name of all available tools
// @Security     ApiKeyAuth
// @Tags         tools
// @Accept       json
// @Produce      json
// @Success      200 {array} []tsf.Spec
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
// @Success      200 {object} tsf.Spec
// @Failure      500 {object} schemas.ErrorResponse
// @Failure      404 {object} schemas.ErrorResponse
// @Router       /v1/tool/{name} [get]
func GetTool(c *fiber.Ctx) error {
	tool, ok := state.Tools[c.Params("name")]
	if !ok {
		return NotFoundError(c)
	}

	return c.JSON(tool)
}

// CreateTool godoc
// @Summary      Create a new spec
// @Description  Upload a new spec to use
// @Security     ApiKeyAuth
// @Tags         tools
// @Accept       json
// @Produce      json
// @Param        spec body tsf.Spec true "new spec"
// @Success      200 {object} tsf.Spec
// @Failure      500 {object} schemas.ErrorResponse
// @Router       /v1/tool [post]
func CreateTool(c *fiber.Ctx) error {
	tool := &tsf.Spec{}
	err := json.Unmarshal(c.Body(), tool)
	if err != nil {
		return BadRequestError(c, []error{err})
	}
	state.Tools[tool.Header.ID] = tool

	return c.JSON(tool)
}
