package main

import (
	"crypto/sha256"
	"log"
	"os"

	_ "git.fuzz.codes/fuzzercloud/workerengine/docs"
	"git.fuzz.codes/fuzzercloud/workerengine/handlers"
	"git.fuzz.codes/fuzzercloud/workerengine/podman"
	"git.fuzz.codes/fuzzercloud/workerengine/state"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

var Mode string

// @title WorkerEngine API
// @version v0.1.0
// @description WorkerEngine is a sandbox API to execute TSF based software.
// @host 127.0.0.1:7001
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes https http
// @BasePath /
func main() {
	var err error

	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	app.Use(logger.New())

	err = godotenv.Load("env")
	if err != nil {
		panic(err)
	}

	Mode = os.Getenv("MODE")

	if Mode == "dev" {
		app.Get("/docs/*", swagger.HandlerDefault)
	} else {
		accessKeyHash = sha256.Sum256([]byte(os.Getenv("ACCESS_KEY")))
		app.Use(keyauth.New(keyauth.Config{
			Validator:    keyValidator,
			ErrorHandler: handlers.UnauthorizedError,
		}))
	}

	err = podman.OpenConnection(os.Getenv("PODMAN_SOCKET_ADDRESS"))
	if err != nil {
		panic(err)
	}

	err = state.ReadTools()
	if err != nil {
		panic(err)
	}

	err = state.ReadTasks()
	if err != nil {
		panic(err)
	}

	RegisterRoutes(app)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
