package task

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	task := app.Group("/task")
	task.Post("/", CreateTask)
	task.Get("/:id", GetTask)
	task.Get("/:id/stdout", GetTaskOutput)
	task.Get("/:id/stream", websocket.New(StreamTaskOutput))
}
