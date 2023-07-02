package main

import (
	"crypto/sha256"
	"log"
	"os"

	_ "fuzz.codes/fuzzercloud/workerengine/docs"
	"fuzz.codes/fuzzercloud/workerengine/handlers"
	"fuzz.codes/fuzzercloud/workerengine/podman"
	"fuzz.codes/fuzzercloud/workerengine/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

var Mode string

// @title WorkerEngine API
// @version 1.0
// @description WorkerEngine is a sandbox API to execute TSF based software.
// @host 127.0.0.1:7001
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	app.Use(logger.New())

	godotenv.Load("env")

	Mode = os.Getenv("MODE")

	if Mode == "dev" {
		app.Get("/docs/*", swagger.HandlerDefault)
	} else {
		accessKeyHash = sha256.Sum256([]byte(os.Getenv("ACCESS_KEY")))
		app.Use(keyauth.New(keyauth.Config{
			Validator: keyValidator,
		}))
	}

	podman.OpenConnection(os.Getenv("PODMAN_SOCKET_ADDRESS"))

	task.ReadTools()

	RegisterRoutes(app)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
