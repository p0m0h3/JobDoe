package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"os"

	_ "fuzz.codes/fuzzercloud/workerengine/docs"
	"fuzz.codes/fuzzercloud/workerengine/task"
	"fuzz.codes/fuzzercloud/workerengine/tool"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

var accessKeyHash [32]byte

func keyValidator(c *fiber.Ctx, k string) (bool, error) {
	hashed := sha256.Sum256([]byte(k))

	if subtle.ConstantTimeCompare(accessKeyHash[:], hashed[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

// @title WorkerEngine API
// @version 1.0
// @description WorkerEngine is a sandbox API to execute TSF based software.
// @host 127.0.0.1:7001
// securitydefinitions.apikey
// @BasePath /
func main() {
	app := fiber.New()

	godotenv.Load("env")

	accessKeyHash = sha256.Sum256([]byte(os.Getenv("ACCESS_KEY")))

	app.Use(keyauth.New(keyauth.Config{
		Validator: keyValidator,
	}))

	app.Get("/docs/*", swagger.HandlerDefault)

	tool.RegisterRoutes(app)
	task.RegisterRoutes(app)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ON")))
}
