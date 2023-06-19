package main

import (
	"log"
	"os"

	_ "fuzz.codes/fuzzercloud/workerengine/docs"
	"fuzz.codes/fuzzercloud/workerengine/tool"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title WorkerEngine API
// @version 1.0
// @description WorkerEngine is a sandbox API to execute TSF based software.
// @host 127.0.0.1:8080
// @BasePath /
func main() {
	app := fiber.New()

	godotenv.Load("env")

	app.Get("/docs/*", swagger.HandlerDefault)

	tool.RegisterRoutes(app)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
