package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"os"
	"systest/log"
	"systest/payloads"
)

func main() {
	logger := log.New("Core")

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "SysTest",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.Error("Internal server error: %s", err.Error())

			return ctx.Status(fiber.StatusInternalServerError).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     "*",
	}))
	app.Use(recover.New())

	RegisterEurekaRoutes(app)
	RegisterApiRoutes(app)

	address := os.Getenv("SERVER_IP_ADDRESS") + ":" + os.Getenv("SERVER_PORT")
	logger.Info("Starting listen server on: %s", address)
	if err := app.Listen(address); err != nil {
		panic(err)
	}
}
