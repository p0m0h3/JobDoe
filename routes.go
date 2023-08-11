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
	task.Delete("/", handlers.PruneTasks)
	task.Get("/:id/log", handlers.GetTaskLog)
	task.Get("/:id/files", handlers.GetTaskOutputFiles)
	task.Get("/:id/stream", websocket.New(handlers.StreamTaskLog))
	task.Get("/:id/stats", handlers.GetTaskStats)
	task.Get("/:id/wait", handlers.WaitOnTask)

	tool := app.Group("/tool")
	tool.Get("/", handlers.GetAllTools)
	tool.Get("/:name", handlers.GetTool)
}
