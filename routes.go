package main

import (
	"git.fuzz.codes/fuzzercloud/workerengine/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterV1Routes(v1 fiber.Router) {
	v1.Get("/", handlers.Ping)

	task := v1.Group("/task")
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

	tool := v1.Group("/tool")
	tool.Get("/", handlers.GetAllTools)
	tool.Get("/:name", handlers.GetTool)
	tool.Put("/", handlers.CreateTool)

}

func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/v1")
	RegisterV1Routes(v1)
}
