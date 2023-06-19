package tool

import (
	"os"
	"strings"

	"fuzz.codes/fuzzercloud/tsf"
	"github.com/gofiber/fiber/v2"
)

// GetAllTools godoc
// @Summary      Get all available tools
// @Description  Get the name of all available tools
// @Tags         tools
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /tool/ [get]
func GetAllTools(c *fiber.Ctx) error {
	toolFiles, err := os.ReadDir(os.Getenv("TOOLS_DIRECTORY"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	result := make([]string, 0)

	for _, file := range toolFiles {
		if strings.HasSuffix(file.Name(), ".toml") {
			result = append(result, strings.TrimSuffix(file.Name(), ".toml"))
		}
	}
	return c.JSON(result)
}

// GetTool godoc
// @Summary      Get available tool by name
// @Description  Get details for a tool
// @Tags         tools
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Failure      404
// @Router       /tool/{name} [get]
func GetTool(c *fiber.Ctx) error {
	data, err := os.ReadFile(os.Getenv("TOOLS_DIRECTORY") + "/" + c.Params("name") + ".toml")
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	tool, err := tsf.Parse(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(tool)
}
