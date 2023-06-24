package task

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	task := app.Group("/task")
	task.Post("/", CreateTask)
	task.Get("/:id", GetTask)

	InitConnection()
}
