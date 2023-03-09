package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
	"systest/log"
	"systest/payloads"
)

var validate validator.Validate

func validatePayload(ctx *fiber.Ctx, payload interface{}) error {
	if err := ctx.BodyParser(payload); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(payloads.NewEurekaError(
			fiber.StatusBadRequest,
			"SysTest could not parse the request body",
			ctx,
		))
		return err
	}

	if err := validate.Struct(payload); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSON(payloads.NewEurekaError(
			fiber.StatusBadRequest,
			"SysTest could not validate the request body",
			ctx,
		))
		return err
	}

	return nil
}

func RegisterEurekaRoutes(app *fiber.App) {
	logger := log.New("Eureka")
	validate = *validator.New()

	baseUrl := os.Getenv("EUREKA_BASE_URL")
	api := app.Group("/eureka")

	logger.Info("Registering Eureka proxy routes...")

	api.Get("/apps/", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for all application instances")
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/apps/delta", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for application instances with deltas")
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/apps/:appId", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for all application instances of: %s", ctx.Params("appId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Post("/apps/:appId", func(ctx *fiber.Ctx) error {
		logger.Info("Registering a new instance of %s", ctx.Params("appId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Delete("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		logger.Info("De-registering %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Put("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		logger.Info("Heartbeat received for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Put("/apps/:appId/:instanceId/status", func(ctx *fiber.Ctx) error {
		logger.Info("Updating status override for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Delete("/apps/:appId/:instanceId/status", func(ctx *fiber.Ctx) error {
		logger.Info("Removing status override for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Put("/apps/:appId/:instanceId/metadata", func(ctx *fiber.Ctx) error {
		logger.Info("Updating metadata for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/instances/:instanceId", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for instance: %s", ctx.Params("instanceId"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/vips/:vipAddress", func(ctx *fiber.Ctx) error {
		logger.Info("Querying instances under vip address: %s", ctx.Params("vipAddress"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})

	api.Get("/svips/:secureVipAddress", func(ctx *fiber.Ctx) error {
		logger.Info("Querying for instances under secure vip address: %s", ctx.Params("secureVipAddress"))
		return proxy.Do(ctx, baseUrl+ctx.OriginalURL())
	})
}
