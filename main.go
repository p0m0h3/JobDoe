package main

import (
	_ "github.com/fuzzercloud/workerengine/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title WorkerEngine API
// @version 1.0
// @description WorkerEngine is a sandbox API to execute TSF based software.
// @host localhost:8080
// @BasePath /
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/docs/*", swagger.HandlerDefault)
	app.Listen(":8080")
}
