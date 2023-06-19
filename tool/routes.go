package tool

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	tool := app.Group("/tool")
	tool.Get("/", GetAllTools)
	tool.Get("/:name", GetTool)

	ReadTools()
}
