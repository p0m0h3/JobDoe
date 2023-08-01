package main

import (
	"fuzz.codes/fuzzercloud/workerengine/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	task := app.Group("/task")
	task.Post("/", handlers.CreateTask)
	task.Get("/", handlers.GetAllTasks)
	task.Get("/:id", handlers.GetTask)
	task.Delete("/:id", handlers.DeleteTask)
	task.Get("/:id/stdout", handlers.GetTaskOutput)
	task.Get("/:id/stream", websocket.New(handlers.StreamTaskOutput))
	task.Get("/:id/stats", handlers.GetTaskStats)
	task.Get("/:id/wait", handlers.WaitOnTask)

	tool := app.Group("/tool")
	tool.Get("/", handlers.GetAllTools)
	tool.Get("/:name", handlers.GetTool)
}
