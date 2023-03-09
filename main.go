package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"os"
	"systest-proxy/log"
	"systest-proxy/payloads"
)

var app *fiber.App

func main() {
	logger := log.New("Core")
	if !fiber.IsChild() {
		logger.Info("Initialising Fiber application...")
	}

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	app = fiber.New(fiber.Config{
		AppName:       "SysTest Proxy",
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			logger.Error("Internal server error: %s", err.Error())

			return ctx.Status(fiber.StatusInternalServerError).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		},
	})

	app.Use(recover.New()) //Panic protection

	if err := app.Listen(":" + os.Getenv("SERVER_PORT")); err != nil {
		panic(err)
	}
}
